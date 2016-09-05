package backends

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bndw/pick/errors"
	"github.com/mitchellh/go-homedir"
)

const (
	safeFileName = "pick.safe"
	safeFileMode = 0600
	safeDirName  = ".pick"
	safeDirMode  = 0700
)

type DiskBackend struct {
	path string
}

func defaultSafePath() *string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	safeDir := fmt.Sprintf("%s/%s", home, safeDirName)

	if _, err := os.Stat(safeDir); os.IsNotExist(err) {
		if mkerr := os.Mkdir(safeDir, safeDirMode); mkerr != nil {
			panic(mkerr)
		}
	}

	safePath := fmt.Sprintf("%s/%s", safeDir, safeFileName)

	return &safePath
}

func NewDiskBackend(path *string) *DiskBackend {
	if path == nil {
		path = defaultSafePath()
	}

	return &DiskBackend{*path}
}

func (db *DiskBackend) Load() ([]byte, error) {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		return nil, &errors.SafeNotFound{}
	}

	return ioutil.ReadFile(db.path)
}

func (db *DiskBackend) Save(data []byte) error {
	return ioutil.WriteFile(db.path, data, safeFileMode)
}
