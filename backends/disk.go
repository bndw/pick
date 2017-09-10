package backends

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

type DiskBackend struct {
	path         string
	backupConfig backupConfig
}

type fileInfoSlice []os.FileInfo

func NewDiskBackend(config Config) (*DiskBackend, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	safePath, ok := config.Settings["path"].(string)
	if ok {
		safePath = formatHomeDir(safePath, homeDir)
	} else {
		safePath, err = defaultSafePath(homeDir)
		if err != nil {
			return nil, err
		}
	}

	config.Backup.DirPath = fmt.Sprintf(defaultBackupDir, homeDir, defaultSafeDirName)

	return &DiskBackend{
		path:         safePath,
		backupConfig: config.Backup,
	}, nil
}

func (db *DiskBackend) Load() ([]byte, error) {
	if _, err := os.Stat(db.path); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrSafeNotFound
		} else {
			return nil, err
		}
	}

	return ioutil.ReadFile(db.path)
}

func (db *DiskBackend) Save(data []byte) error {
	tmpFile := db.path + ".tmp"
	if err := ioutil.WriteFile(tmpFile, data, defaultSafeFileMode); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}
	if err := os.Rename(tmpFile, db.path); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}
	return nil
}

func (db *DiskBackend) SafeLocation() string {
	return db.path
}

func (f fileInfoSlice) Len() int {
	return len(f)
}

func (f fileInfoSlice) Less(i, j int) bool {
	return f[i].ModTime().Before(f[j].ModTime())
}

func (f fileInfoSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (db *DiskBackend) cleanOldBackups(max int) error {
	files, err := ioutil.ReadDir(db.backupConfig.DirPath)
	if err != nil {
		return err
	}

	filesSorted := make(fileInfoSlice, 0, len(files))
	for _, f := range files {
		filesSorted = append(filesSorted, f)
	}
	sort.Sort(filesSorted)
	max = min(max, len(filesSorted))

	for _, f := range filesSorted[:len(filesSorted)-max] {
		p := fmt.Sprintf("%s/%s", db.backupConfig.DirPath, f.Name())
		if err := os.Remove(p); err != nil {
			fmt.Println("Error removing old backup", err.Error())
		}
	}

	return nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (db *DiskBackend) Backup() error {
	if db.backupConfig.MaxFiles == 0 {
		// Keep no backups
		_ = db.cleanOldBackups(0)
		return errors.ErrBackupDisabled
	} else if db.backupConfig.MaxFiles > 0 {
		// Subtract one as we are about to create another backup
		if err := db.cleanOldBackups(db.backupConfig.MaxFiles - 1); err != nil {
			if !os.IsNotExist(err) {
				fmt.Println("Failed to remove old backup(s)", err.Error())
			}
		}
	}

	backupDir := db.backupConfig.DirPath
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
		return errors.ErrBackupFileExists
	}

	data, err := db.Load()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(backupPath, data, defaultSafeFileMode)
}

func defaultSafePath(homeDir string) (string, error) {
	safeDir := fmt.Sprintf("%s/%s", homeDir, defaultSafeDirName)

	if _, err := os.Stat(safeDir); err != nil {
		if os.IsNotExist(err) {
			if mkerr := os.Mkdir(safeDir, defaultSafeDirMode); mkerr != nil {
				return "", mkerr
			}
		} else {
			return "", err
		}
	}

	safePath := fmt.Sprintf("%s/%s", safeDir, defaultSafeFileName)

	return safePath, nil
}

func formatHomeDir(str, home string) string {
	return strings.Replace(str, "$HOME", home, 1)
}
