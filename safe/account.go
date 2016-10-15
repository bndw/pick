package safe

import (
	"time"
)

type Account struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedOn int64  `json:"createdOn"`
}

func NewAccount(name, username, password string) *Account {
	return &Account{
		Username:  username,
		Password:  password,
		CreatedOn: time.Now().Unix(),
	}
}
