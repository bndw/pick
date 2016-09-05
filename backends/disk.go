package backends

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	homeDir  string
)

type DiskBackend struct {
	path string
}

func NewDiskBackend(config Config) (*DiskBackend, error) {
	var err error
	if homeDir, err = homedir.Dir(); err != nil {
		return nil, err
	}

	safePath, ok := config.Settings["path"].(string)
	if ok && strings.HasPrefix(safePath, "$HOME") {
		safePath = formatHomeDir(safePath, homeDir)
	} else {
		safePath, err = defaultSafePath()
		if err != nil {
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
	safeDir := fmt.Sprintf("%s/%s", homeDir, defaultSafeDirName)

	if _, err := os.Stat(safeDir); os.IsNotExist(err) {
		if mkerr := os.Mkdir(safeDir, defaultSafeDirMode); mkerr != nil {
			return "", mkerr
		}
	}

	safePath := fmt.Sprintf("%s/%s", safeDir, defaultSafeFileName)

	return safePath, nil
}

func formatHomeDir(str, home string) string {
	return strings.Replace(str, "$HOME", home, 1)
}
