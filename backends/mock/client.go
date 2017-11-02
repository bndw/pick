package mock

// client is used for tests
type client struct {
	Data []byte
}

func (c *client) Load() ([]byte, error) {
	return c.Data, nil
}

func (c *client) Save(ciphertext []byte) error {
	c.Data = ciphertext
	return nil
}

func (c *client) Backup() error {
	return nil
}

func (c *client) SafeLocation() string {
	return "mock-safe-location"
}
