package safe

import (
	"github.com/bndw/pick/errors"
)

func (s *Safe) Add(name, username, password string) error {
	if _, exists := s.Accounts[name]; exists {
		return &errors.AccountExists{}
	}

	s.Accounts[name] = NewAccount(name, username, password)

	return s.save()
}
