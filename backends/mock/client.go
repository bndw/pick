package mock

import (
	"github.com/bndw/pick/errors"
)

// client is used for tests
type client struct {
	writable bool

	Data []byte
}

func (c *client) Load() ([]byte, error) {
	return c.Data, nil
}

func (c *client) Save(ciphertext []byte) error {
	if !c.writable {
		return errors.ErrSafeNotWritable
	}

	c.Data = ciphertext
	return nil
}

func (c *client) Backup() error {
	return nil
}

func (c *client) SafeLocation() string {
	return "mock-safe-location"
}

func (c *client) IsWritable() bool {
	return c.writable
}

func (c *client) SetWritable(writable bool) error {
	c.writable = writable
	return nil
}

func (c *client) Lock() error {
	return nil
}

func (c *client) Unlock() error {
	return nil
}
