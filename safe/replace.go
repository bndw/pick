package safe

func (s *Safe) Replace(name, username, password string) (*Account, error) {
	if _, err := s.Get(name); err != nil {
		return nil, err
	}

	account := NewAccount(name, username, password)
	s.Accounts[name] = *account

	if err := s.save(); err != nil {
		return nil, err
	}

	return account, nil
}
