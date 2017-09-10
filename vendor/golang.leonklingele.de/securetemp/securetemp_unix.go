// +build !darwin

package securetemp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	devShmPath = "/dev/shm"
)

func createRAMDisk(size int) (string, func(), error) {
	var (
		mountPath string
		err       error

		doCleanup = true
	)
	cleanupFunc := func() { cleanupRAMDisk(mountPath) }
	defer func() {
		if doCleanup {
			cleanupFunc()
		}
	}()

	if unix.Access(devShmPath, unix.W_OK) == nil && unix.Access(devShmPath, unix.X_OK) == nil {
		// We'll use /dev/shm
		mountPath, err = ioutil.TempDir(devShmPath, globalPrefix)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create temp dir in %s: %s", devShmPath, err)
		}
		doCleanup = false
		return mountPath, cleanupFunc, nil
	}

	mountPath, err = ioutil.TempDir("", globalPrefix)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %s", err)
	}

	if err := createAndMount(mountPath, size); err != nil {
		return "", nil, err
	}

	doCleanup = false
	return mountPath, cleanupFunc, nil
}

func cleanupRAMDisk(mountPath string) {
	if mountPath != "" {
		if !strings.HasPrefix(mountPath, devShmPath) {
			exec.Command("umount", mountPath).Run()
		}
		os.RemoveAll(mountPath)
	}
}

func createAndMount(mountPath string, size int) error {
	path, err := exec.LookPath("mount_mfs")
	if err != nil {
		return errors.New("did not find 'mount_mfs'")
	}

	cmd := exec.Command(path,
		"-o", "noatime",
		fmt.Sprintf("-s %d", size*1024*2),
		"/dev/wd0b",
		mountPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'mount_mfs': %s", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("'mount_mfs' failed: %s, %s", err, stderr.String())
	}

	return nil
}
