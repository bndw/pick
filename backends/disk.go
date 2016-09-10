package backends

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/bndw/pick/errors"
	"github.com/mitchellh/go-homedir"
)

const (
	defaultSafeFileMode     = 0600
	defaultSafeFileName     = "pick.safe"
	defaultSafeDirMode      = 0700
	defaultSafeDirName      = ".pick"
	defaultBackupDir        = "%s/%s/backups"
	defaultBackupFileName   = "pick_%s.safe"
	defaultBackupTimeFormat = "2006-01-02_15-04-05"
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
	if ok {
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

func (db *DiskBackend) Backup() error {
	backupDir := fmt.Sprintf(defaultBackupDir, homeDir, defaultSafeDirName)
	timeFormat := time.Now().Format(defaultBackupTimeFormat)
	backupFileName := fmt.Sprintf(defaultBackupFileName, timeFormat)
	backupPath := backupDir + "/" + backupFileName

	if _, err := os.Stat(backupDir); err != nil {
		if os.IsNotExist(err) {
			if mkerr := os.Mkdir(backupDir, defaultSafeDirMode); mkerr != nil {
				return mkerr
			}
		} else {
			return err
		}
	}

	if _, err := os.Stat(backupPath); err == nil {
		return &errors.BackupFileExists{}
	}

	data, err := db.Load()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(backupPath, data, defaultSafeFileMode)
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
