package models

type LoginDto struct {
	AccountId string `json:"accountId"`
	Password  string `json:"password"`
}
