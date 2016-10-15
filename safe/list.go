package safe

func (s *Safe) List() map[string]Account {
	return s.Accounts
}
