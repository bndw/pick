package backends

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bndw/pick/errors"
	"github.com/mitchellh/go-homedir"
)

const (
	defaultSafeFileMode = 0600
	defaultSafeFileName = "pick.safe"
	defaultSafeDirMode  = 0700
	defaultSafeDirName  = ".pick"
)

var (
	safePath string
)

type DiskBackend struct {
	path string
}

func NewDiskBackend(config Config) (*DiskBackend, error) {
	safePath, ok := config.Settings["path"].(string)
	if !ok {
		// Use default path when not specified
		var err error
		if safePath, err = defaultSafePath(); err != nil {
			return nil, err
		}
	}

	return &DiskBackend{safePath}, nil
}

func (db *DiskBackend) Load() ([]byte, error) {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		return nil, &errors.SafeNotFound{}
	}

	return ioutil.ReadFile(db.path)
}

func (db *DiskBackend) Save(data []byte) error {
	return ioutil.WriteFile(db.path, data, defaultSafeFileMode)
}

func defaultSafePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	safeDir := fmt.Sprintf("%s/%s", home, defaultSafeDirName)

	if _, err := os.Stat(safeDir); os.IsNotExist(err) {
		if mkerr := os.Mkdir(safeDir, defaultSafeDirMode); mkerr != nil {
			return "", mkerr
		}
	}

	safePath := fmt.Sprintf("%s/%s", safeDir, defaultSafeFileName)

	return safePath, nil
}
