package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/bndw/pick/errors"
)

const (
	defaultBackupDir        = "%s/%s/backups"
	defaultBackupFileName   = "pick_%s.safe"
	defaultBackupTimeFormat = "2006-01-02_15-04-05"
)

type fileInfoSlice []os.FileInfo

func (f fileInfoSlice) Len() int {
	return len(f)
}

func (f fileInfoSlice) Less(i, j int) bool {
	return f[i].ModTime().Before(f[j].ModTime())
}

func (f fileInfoSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (c *client) cleanOldBackups(max int) error {
	files, err := ioutil.ReadDir(c.backupConfig.DirPath)
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
		p := fmt.Sprintf("%s/%s", c.backupConfig.DirPath, f.Name())
		if err := os.Remove(p); err != nil {
			fmt.Println("Error removing old backup", err.Error())
		}
	}

	return nil
}

func (c *client) Backup() error {
	if c.backupConfig.MaxFiles == 0 {
		// Keep no backups
		_ = c.cleanOldBackups(0)
		return errors.ErrBackupDisabled
	} else if c.backupConfig.MaxFiles > 0 {
		// Subtract one as we are about to create another backup
		if err := c.cleanOldBackups(c.backupConfig.MaxFiles - 1); err != nil {
			if !os.IsNotExist(err) {
				fmt.Println("Failed to remove old backup(s)", err.Error())
			}
		}
	}

	backupDir := c.backupConfig.DirPath
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

	data, err := c.Load()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(backupPath, data, defaultSafeFileMode)
}
