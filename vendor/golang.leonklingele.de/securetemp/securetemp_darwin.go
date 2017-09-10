package securetemp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func createRAMDisk(size int) (string, func(), error) {
	var (
		mountPath, devPath string
		err                error

		doCleanup = true
	)
	cleanupFunc := func() { cleanupRAMDisk(devPath, mountPath) }
	defer func() {
		if doCleanup {
			cleanupFunc()
		}
	}()

	mountPath, err = ioutil.TempDir("", globalPrefix)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %s", err)
	}

	devPath, err = createDev(size)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create device: %s", err)
	}

	if err := formatDev(devPath); err != nil {
		return "", nil, fmt.Errorf("failed to format device: %s", err)
	}
	if err := mountDev(devPath, mountPath); err != nil {
		return "", nil, fmt.Errorf("failed to mount device: %s", err)
	}

	doCleanup = false
	return mountPath, cleanupFunc, nil
}

func cleanupRAMDisk(devPath, mountPath string) {
	if mountPath != "" {
		exec.Command("umount", mountPath).Run()
		defer func() {
			os.RemoveAll(mountPath)
		}()
	}
	if devPath != "" {
		exec.Command("diskutil", "quiet", "eject", devPath).Run()
	}
}

func createDev(size int) (string, error) {
	path, err := exec.LookPath("hdid")
	if err != nil {
		return "", errors.New("did not find 'hdid'")
	}

	cmd := exec.Command(path,
		"-drivekey",
		"system-image=yes",
		"-nomount",
		fmt.Sprintf("ram://%d", size*1024*2))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start 'hdid': %s", err)
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("'hdid' failed: %s, %s", err, stderr.String())
	}

	return strings.Split(stdout.String(), " ")[0], nil
}

func formatDev(devPath string) error {
	path, err := exec.LookPath("newfs_hfs")
	if err != nil {
		return errors.New("did not find 'newfs_hfs'")
	}

	cmd := exec.Command(path,
		"-M", "700",
		devPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'newfs_hfs': %s", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("'newfs_hfs' failed: %s, %s", err, stderr.String())
	}

	return nil
}

func mountDev(devPath, mountPath string) error {
	path, err := exec.LookPath("mount")
	if err != nil {
		return errors.New("did not find 'mount'")
	}

	cmd := exec.Command(path,
		"-t", "hfs",
		"-o", "noatime",
		"-o", "nobrowse",
		devPath, mountPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'mount': %s", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("'mount' failed: %s, %s", err, stderr.String())
	}

	return nil
}
