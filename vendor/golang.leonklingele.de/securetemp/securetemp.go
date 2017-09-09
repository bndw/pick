package securetemp

import (
	"io/ioutil"
	"os"
)

const (
	// SizeKB represents one Kilobyte
	SizeKB = 1 << (10 * iota)
	// SizeMB represents one Megabyte (1000 KB)
	SizeMB
	// SizeGB represents one Gigabyte (1000 MB)
	SizeGB
)

const (
	// DefaultSize specifies the default size for a RAM disk in MB
	DefaultSize = 4 * SizeMB

	// globalPrefix will be used if we need a prefix for
	// something (e.g. files & folders)
	globalPrefix = "securetemp"
)

// TempDir creates a new RAM disk with size 'size' (in bytes)
// and returns the path to it.
// Use this function only if you intend to create multiple
// files inside your RAM disk, else prefer to use 'TempFile'.
func TempDir(size int) (string, func(), error) {
	path, cleanupFunc, err := createRAMDisk(size)
	if err != nil {
		return "", nil, err
	}
	return path, cleanupFunc, nil
}

// TempFile creates a new RAM disk with size 'size' (in bytes),
// creates a temp file in it and returns a pointer to that file.
// Use this function only if you intend to create a single
// inside your RAM disk, else prefer to use 'TempDir'.
func TempFile(size int) (*os.File, func(), error) {
	path, cleanupFunc, err := TempDir(size)
	if err != nil {
		return nil, nil, err
	}
	file, err := ioutil.TempFile(path, globalPrefix)
	if err != nil {
		return nil, nil, err
	}
	return file, func() { file.Close(); cleanupFunc() }, nil
}
