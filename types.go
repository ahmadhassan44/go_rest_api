package main

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

type Account struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(math.MaxInt64)),
	}
}
