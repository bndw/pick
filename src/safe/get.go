package safe

import (
	"github.com/bndw/pick/errors"
)

func (s *Safe) Get(name string) (*Account, error) {
	account, exists := s.Accounts[name]
	if !exists {
		return nil, &errors.AccountNotFound{}
	}

	return &account, nil
}
