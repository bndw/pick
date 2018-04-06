// +build windows

package single

import (
	"fmt"
	"os"
	"path/filepath"
)

// Filename returns an absolute filename, appropriate for the operating system
func (s *Single) Filename() string {
	if len(Lockfile) > 0 {
		return Lockfile
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s.lock", s.name))
}

// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s *Single) CheckLock() error {

	if err := os.Remove(s.Filename()); err != nil && !os.IsNotExist(err) {
		return ErrAlreadyRunning
	}

	file, err := os.OpenFile(s.Filename(), os.O_EXCL|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.file = file

	return nil
}

// TryUnlock closes and removes the lockfile
func (s *Single) TryUnlock() error {
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close the lock file: %v", err)
	}
	if err := os.Remove(s.Filename()); err != nil {
		return fmt.Errorf("failed to remove the lock file: %v", err)
	}
	return nil
}
