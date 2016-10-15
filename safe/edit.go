package safe

func (s *Safe) Edit(name, username, password string) (*Account, error) {
	oldAcc, err := s.Get(name)
	if err != nil {
		return nil, err
	}

	// Make a cheap copy of the current account
	newAcc := *oldAcc
	// Update the new account
	newAcc.Update(func(acc *Account) {
		if username != "" {
			acc.Username = username
		}
		if password != "" {
			acc.Password = password
		}
	})
	// Remove history of the old account as the new account already has it
	oldAcc.History = nil
	// Add the old account to the new account's history
	newAcc.History = append(newAcc.History, *oldAcc)
	// Store new account
	s.Accounts[name] = newAcc

	if err := s.save(); err != nil {
		return nil, err
	}

	return &newAcc, nil
}
