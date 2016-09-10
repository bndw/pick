package safe

import (
	"github.com/bndw/pick/errors"
)

func (s *Safe) Replace(name, username, password string) error {
	if _, exists := s.Accounts[name]; exists {
		s.Accounts[name] = NewAccount(name, username, password)
		return s.save()
	}

	return &errors.AccountNotFound{}
}
