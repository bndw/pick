package safe

import (
	"github.com/bndw/pick/errors"
)

func (s *Safe) Add(name, username, password string) (*Account, error) {
	if existingAccount, err := s.Get(name); err == nil {
		return existingAccount, errors.ErrAccountAlreadyExists
	}

	account := NewAccount(username, password)
	s.Accounts[name] = *account

	if err := s.save(); err != nil {
		return nil, err
	}

	return account, nil
}
