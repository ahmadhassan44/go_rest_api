package models

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	UserName       string `json:"userName"`
	HashedPassword string
	Number         int64     `json:"number"`
	Balance        int64     `json:"balance"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewAccount(firstName string, lastName string, userName string, hashedPassword string) *Account {
	return &Account{
		ID:             uuid.NewString(),
		FirstName:      firstName,
		LastName:       lastName,
		UserName:       userName,
		HashedPassword: hashedPassword,
		Number:         int64(rand.Intn(math.MaxInt64)),
	}
}

type CreateAccountDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
}
type UpdateAccountDto struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Balance   *int64  `json:"balance"`
}
type TransferMoneyDto struct {
	ReceiverId string `json:"receiverId"`
	Amount     int64  `json:"amount"`
}
type AccountError struct {
	Msg        string
	StatusCode int
}

func (e AccountError) Error() string {
	return e.Msg
}
func NewAccountError(message string, statusCode int) *AccountError {
	return &AccountError{
		Msg:        message,
		StatusCode: statusCode,
	}
}
