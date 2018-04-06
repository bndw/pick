package file

import (
	"io/ioutil"
	"os"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/errors"
	"github.com/marcsauter/single"
)

const (
	defaultSafeFileMode = 0600
	defaultSafeFileName = "pick.safe"

	defaultSafeDirMode = 0700
	defaultSafeDirName = ".pick"
)

type client struct {
	writable bool
	lock     *single.Single

	path         string
	backupConfig backends.BackupConfig
}

func (c *client) Load() ([]byte, error) {
	if _, err := os.Stat(c.path); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrSafeNotFound
		} else {
			return nil, err
		}
	}

	return ioutil.ReadFile(c.path)
}

func (c *client) Save(data []byte) error {
	if !c.writable {
		return errors.ErrSafeNotWritable
	}

	tmpFile := c.path + ".tmp"
	if err := ioutil.WriteFile(tmpFile, data, defaultSafeFileMode); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}
	if err := os.Rename(tmpFile, c.path); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}
	return nil
}

func (c *client) SafeLocation() string {
	return c.path
}

func (c *client) IsWritable() bool {
	return c.writable
}

func (c *client) SetWritable(writable bool) error {
	c.writable = writable
	if writable {
		// NOTE(leon): No need to call c.Unlock(). The lock is automatically released again once this instance of pick terminates.
		if err := c.Lock(); err == single.ErrAlreadyRunning {
			return errors.ErrAlreadyRunning
		} else {
			return err
		}
	}
	return nil
}

func (c *client) Lock() error {
	return c.lock.CheckLock()
}

func (c *client) Unlock() error {
	return c.lock.TryUnlock()
}
