package models

type LoginDto struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
