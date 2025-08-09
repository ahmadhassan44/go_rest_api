package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(*Account) error
	DeleteAccount(*Account) error
	GetAccountById(string) (*Account, error)
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
	query := `CREATE TABLE IF  NOT EXISTS account(
		id varchar(36) primary key,
		first_name varchar(100),
		last_name varchar(100),
		number int,
		balance int,
		created_at timestamp,
		updated_at timestamp
	)`
	_, err := pgStore.db.Exec(query)
	return err
}

func (pgStore *PostGresStore) CreateAccount(*Account) error {
	return nil

}
func (pgStore *PostGresStore) DeleteAccount(*Account) error {
	return nil

}
func (pgStore *PostGresStore) UpdateAccount(*Account) error {
	return nil

}
func (pgStore *PostGresStore) GetAccountById(id string) (*Account, error) {
	return nil, nil

}
