package safe

func (s *Safe) List() []Account {
	var accounts []Account
	for _, account := range s.Accounts {
		accounts = append(accounts, account)
	}

	return accounts
}
