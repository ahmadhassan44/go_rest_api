package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(string) error
	GetAccountById(string) (*Account, error)
	GetAllAccounts() ([]*Account, error)
	UpdateAccount(*Account) error
}
type Storage interface {
	AccountStorage
}

type PostGresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostGresStore, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=gobank password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostGresStore{db: db}, nil
}

func (pgStore *PostGresStore) Init() error {
	return pgStore.createAccountTable()
}

func (pgStore *PostGresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account(
		id uuid primary key,
		first_name varchar(100),
		last_name varchar(100),
		number bigint,
		balance int,
		created_at timestamp default now(),
		updated_at timestamp default now()
	)`
	_, err := pgStore.db.Exec(query)
	return err
}

func (pgStore *PostGresStore) CreateAccount(account *Account) (*Account, error) {
	query := `INSERT INTO 
	account(id,first_name,last_name,number,balance) 
	VALUES ($1,$2,$3,$4,0)`
	_, err := pgStore.db.Exec(query, account.ID, account.FirstName, account.LastName, account.Number)
	if err != nil {
		return nil, err
	}

	var createdAccount *Account
	createdAccount, err = pgStore.GetAccountById(account.ID)

	if err != nil {
		return nil, err
	}

	return createdAccount, nil
}
func (pgStore *PostGresStore) DeleteAccount(id string) error {
	query := "DELETE FROM account WHERE id= $1"
	res, err := pgStore.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return NewAccountError(
			fmt.Sprintf("Account with ID: %s not found!", id), http.StatusNotFound,
		)
	}
	return nil
}
func (pgStore *PostGresStore) UpdateAccount(*Account) error {
	return nil

}
func (pgStore *PostGresStore) GetAccountById(id string) (*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at, updated_at 
	FROM account WHERE id = $1`
	var account Account
	row := pgStore.db.QueryRow(query, id)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAccountError(
				fmt.Sprintf("Account with ID: %s not found!", id), http.StatusNotFound,
			)
		}
		return nil, err
	}
	return &account, nil
}
func (pgStore *PostGresStore) GetAllAccounts() ([]*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at, updated_at 
	FROM account`
	rows, err := pgStore.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
