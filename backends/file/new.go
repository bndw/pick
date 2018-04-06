package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/bndw/pick/backends"
	"github.com/marcsauter/single"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	lockFileName = "pick"
)

func _new(config *backends.Config) (backends.Client, error) {
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

	lock := single.New(lockFileName)
	c := &client{
		lock:         lock,
		path:         safePath,
		backupConfig: config.Backup,
	}

	return c, nil
}

func formatHomeDir(str, home string) string {
	return strings.Replace(str, "$HOME", home, 1)
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
