package securetemp

import (
	"errors"
)

func createRAMDisk(size int) (string, func(), error) {
	return mountPath, func() { return errors.New("windows is currently unsupported") }, nil
}
