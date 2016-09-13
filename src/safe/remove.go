package safe

import (
	"fmt"
)

func (s *Safe) Remove(name string) error {
	_, exists := s.Accounts[name]
	if !exists {
		return fmt.Errorf("Account not found")
	}

	delete(s.Accounts, name)

	return s.save()
}
