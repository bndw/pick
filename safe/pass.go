package safe

func (s *Safe) ChangePassword(p []byte) error {
	s.password = p
	return s.save()
}
