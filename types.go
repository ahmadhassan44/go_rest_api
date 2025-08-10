package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(math.MaxInt64)),
	}
}

type CreateAccountDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
