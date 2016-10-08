package safe

import (
	"time"
)

type Account struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedOn int64  `json:"createdOn"`
}

func NewAccount(name, username, password string) *Account {
	return &Account{
		Name:      name,
		Username:  username,
		Password:  password,
		CreatedOn: time.Now().Unix(),
	}
}
