package safe

func (s *Safe) Init() error {
	return s.save()
}
