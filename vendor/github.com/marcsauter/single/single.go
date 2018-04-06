// package single provides a mechanism to ensure, that only one instance of a program is running

package single

import (
	"errors"
	"log"
	"os"
)

var (
	// ErrAlreadyRunning -- the instance is already running
	ErrAlreadyRunning = errors.New("the program is already running")
	// Lockfile -- the lock file to check
	Lockfile string
)

// Single represents the name and the open file descriptor
type Single struct {
	name string
	file *os.File
}

// New creates a Single instance
func New(name string) *Single {
	return &Single{name: name}
}

// Lock tries to obtain an exclude lock on a lockfile and exits the program if an error occurs
func (s *Single) Lock() {
	if err := s.CheckLock(); err != nil {
		log.Fatal(err)
	}
}

// Unlock releases the lock, closes and removes the lockfile. All errors will be reported directly.
func (s *Single) Unlock() {
	if err := s.TryUnlock(); err != nil {
		log.Print(err)
	}
}
