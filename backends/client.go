package backends

type Client interface {
	Load() (data []byte, err error)
	Save(data []byte) error
	Backup() error
	SafeLocation() string
	IsWritable() bool
	SetWritable(writable bool) error
	Lock() error
	Unlock() error
}

func New(config *Config) (Client, error) {
	c, err := clientByName(config.Type)
	if err != nil {
		return nil, err
	}
	return c.newClientFunc(config)
}

func NewWithType(t string, config *Config) (Client, error) {
	config.Type = t
	return New(config)
}
