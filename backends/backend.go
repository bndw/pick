package backends

type Backend interface {
	Load() ([]byte, error)
	Save([]byte) error
}
