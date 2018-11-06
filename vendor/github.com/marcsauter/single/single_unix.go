// +build linux solaris darwin freebsd openbsd netbsd dragonfly

package single

import (
	"fmt"
	"os"
	"syscall"
)

// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s *Single) CheckLock() error {

	// open/create lock file
	f, err := os.OpenFile(s.Filename(), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.file = f
	// set the lock type to F_WRLCK, therefor the file has to be opened writable
	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	}
	// try to obtain an exclusive lock - FcntlFlock seems to be the portable *ix way
	if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return ErrAlreadyRunning
	}

	return nil
}

// TryUnlock unlocks, closes and removes the lockfile
func (s *Single) TryUnlock() error {
	// set the lock type to F_UNLCK
	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
		Pid:  int32(os.Getpid()),
	}
	if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return fmt.Errorf("failed to unlock the lock file: %v", err)
	}
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close the lock file: %v", err)
	}
	if err := os.Remove(s.Filename()); err != nil {
		return fmt.Errorf("failed to remove the lock file: %v", err)
	}
	return nil
}
