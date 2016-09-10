package backends

type Client interface {
	Load() ([]byte, error)
	Save([]byte) error
	Backup() error
}

func New(config *Config) (Client, error) {
	switch config.Type {
	default:
		fallthrough
	case "file":
		return NewDiskBackend(*config)
	case "mock":
		return NewMockBackend(), nil
	}
}
