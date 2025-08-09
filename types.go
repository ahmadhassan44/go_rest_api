package main

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	FirstName string
	LastName  string
	Number    int64
	Balance   int64
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(math.MaxInt64)),
	}
}
