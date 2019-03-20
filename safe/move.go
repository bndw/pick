package safe

import (
	"fmt"
)

func (s *Safe) Move(name, newName string) error {
	account, err := s.Get(name)
	if err != nil {
		return err
	}

	if name == newName {
		return fmt.Errorf("New name must be different")
	}
	if _, err := s.Get(newName); err == nil {
		return fmt.Errorf("Account with new name already exists")
	}

	s.Accounts[newName] = *account
	delete(s.Accounts, name)

	return s.save()
}
