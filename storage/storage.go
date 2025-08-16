package storage

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	models "github.com/ahmadhassan44/go_rest_api/models"
	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(*models.Account) (*models.Account, error)
	DeleteAccount(string) error
	GetAccountById(string) (*models.Account, error)
	GetAllAccounts() ([]*models.Account, error)
	UpdateAccount(string, *models.UpdateAccountDto) error
	TransferMoney(*models.TransferMoneyDto) error
	GetUserByUserName(string) (*string, error)
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

	err := pgStore.createAccountTable()
	if err != nil {
		return err
	}
	err = pgStore.createRefreshTokensTable()
	return err
}
func (pgStore *PostGresStore) createRefreshTokensTable() error {
	query := `CREATE TABLE IF NOT EXISTS refresh_tokens (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE,
		token VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		revoked BOOLEAN NOT NULL DEFAULT FALSE
	)`
	_, err := pgStore.db.Exec(query)
	if err != nil {
		return err
	}
	query = `CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);`
	_, err = pgStore.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (pgStore *PostGresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account(
		id uuid primary key,
		first_name varchar(100),
		last_name varchar(100),
		user_name varchar(100),
		password_hash varchar(100),
		number bigint,
		balance int,
		created_at timestamp default now(),
		updated_at timestamp default now()
	)`
	_, err := pgStore.db.Exec(query)
	return err
}

func (pgStore *PostGresStore) CreateAccount(account *models.Account) (*models.Account, error) {
	query := `INSERT INTO 
	account(id,first_name,last_name,user_name,password_hash,number,balance) 
	VALUES ($1,$2,$3,$4,$5,$6,0)`
	_, err := pgStore.db.Exec(query, account.ID, account.FirstName, account.LastName, account.UserName, account.HashedPassword, account.Number)
	if err != nil {
		return nil, err
	}

	var createdAccount *models.Account
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
		return models.NewAccountError(
			fmt.Sprintf("Account with ID: %s not found!", id), http.StatusNotFound,
		)
	}
	return nil
}
func (pgStore *PostGresStore) UpdateAccount(id string, updateAccountDto *models.UpdateAccountDto) error {
	updates := map[string]any{}
	if updateAccountDto.FirstName != nil {
		updates["first_name"] = *updateAccountDto.FirstName
	}
	if updateAccountDto.LastName != nil {
		updates["last_name"] = *updateAccountDto.LastName
	}
	if updateAccountDto.Balance != nil {
		updates["balance"] = *updateAccountDto.Balance
	}
	if len(updates) == 0 {
		return models.NewAccountError("Nothing specified to update!", http.StatusBadRequest)
	}
	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	for col, val := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, argPos))
		args = append(args, val)
		argPos++
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE account SET %s WHERE id = $%d",
		strings.Join(setParts, ", "),
		argPos,
	)
	res, err := pgStore.db.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.NewAccountError(
			fmt.Sprintf("Account with ID: %s not found!", id), http.StatusNotFound,
		)
	}
	return nil
}
func (pgStore *PostGresStore) GetAccountById(id string) (*models.Account, error) {
	query := `SELECT id, first_name, last_name,user_name, number, balance, created_at, updated_at 
	FROM account WHERE id = $1`
	var account models.Account
	row := pgStore.db.QueryRow(query, id)
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.UserName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewAccountError(
				fmt.Sprintf("Account with ID: %s not found!", id), http.StatusNotFound,
			)
		}
		return nil, err
	}
	return &account, nil
}
func (pgStore *PostGresStore) GetAllAccounts() ([]*models.Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at, updated_at 
	FROM account`
	rows, err := pgStore.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*models.Account{}
	for rows.Next() {
		account := &models.Account{}
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
func (pgStore *PostGresStore) TransferMoney(transferMoneyDto *models.TransferMoneyDto) error {
	tx, err := pgStore.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// TODO: replace with sender id coming from jwt
	senderID, receiverID := "94e531b7-5c08-4297-a889-77298034bc32", transferMoneyDto.ReceiverId
	var a, b string
	if senderID < receiverID {
		a, b = senderID, receiverID
	} else {
		a, b = receiverID, senderID
	}
	//acquire lock
	if _, err := tx.Exec(`SELECT id FROM account WHERE id IN ($1, $2) FOR UPDATE`, a, b); err != nil {
		return err
	}
	//debit
	res, err := tx.Exec(
		`UPDATE account SET balance = balance - $1, updated_at = now()
		 WHERE id = $2 AND balance >= $1`, transferMoneyDto.Amount, senderID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return models.NewAccountError("insufficient funds", http.StatusBadRequest)
	}
	// Credit
	if _, err := tx.Exec(
		`UPDATE account SET balance = balance + $1, updated_at = now()
		 WHERE id = $2`, transferMoneyDto.Amount, receiverID,
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (pgStore *PostGresStore) GetUserByUserName(userName string) (*string, error) {
	query := "SELECT user_name,password_hash FROM account WHERE user_name = $1"
	row := pgStore.db.QueryRow(query, userName)
	var dbUser string
	var hashedPassword string
	err := row.Scan(
		&dbUser,
		&hashedPassword,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewAccountError(
				fmt.Sprintf("Account with username: %s not found!", userName), http.StatusNotFound,
			)
		}
		return nil, err
	}
	return &hashedPassword, nil
}
