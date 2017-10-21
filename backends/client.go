package backends

type Client interface {
	Load() (data []byte, err error)
	Save(data []byte) error
	Backup() error
	SafeLocation() string
}

func New(config *Config) (Client, error) {
	switch config.Type {
	default:
		fallthrough
	case ConfigTypeFile:
		return NewDiskBackend(*config)
	case ConfigTypeS3:
		return NewS3Backend(*config)
	case ConfigTypeMock:
		return NewMockBackend(), nil
	}
}
