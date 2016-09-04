package safe

import (
	"github.com/atotto/clipboard"
	"github.com/bndw/pick/errors"
)

func (s *Safe) Copy(name string) error {
	account, exists := s.Accounts[name]
	if !exists {
		return &errors.AccountNotFound{}
	}

	return clipboard.WriteAll(account.Password)
}
