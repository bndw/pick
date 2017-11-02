package file

import (
	"io/ioutil"
	"os"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/errors"
)

const (
	defaultSafeFileMode = 0600
	defaultSafeFileName = "pick.safe"

	defaultSafeDirMode = 0700
	defaultSafeDirName = ".pick"
)

type client struct {
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
