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

func NewPostgresStore() (*PostGresStore, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	connStr := fmt.Sprintf("host=%s user=%s dbname=gobank password=%s sslmode=disable", dbHost, dbUser, dbPass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostGresStore{db: db}, nil

}
