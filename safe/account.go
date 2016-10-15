package safe

import (
	"time"
)

type Account struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	CreatedOn  int64     `json:"createdOn"`
	ModifiedOn int64     `json:"modifiedOn"`
	History    []Account `json:"history,omitempty"`
}

func (acc *Account) UpdateLastModifed() {
	acc.ModifiedOn = time.Now().Unix()
}

func NewAccount(name, username, password string) *Account {
	ts := time.Now().Unix()
	return &Account{
		Username:   username,
		Password:   password,
		CreatedOn:  ts,
		ModifiedOn: ts,
		History:    make([]Account, 0),
	}
}
