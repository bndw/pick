package single

import (
	"fmt"
	"path/filepath"
)

// Filename returns an absolute filename, appropriate for the operating system
func (s *Single) Filename() string {
	if len(Lockfile) > 0 {
		return Lockfile
	}
	return filepath.Join("/var/lock", fmt.Sprintf("%s.lock", s.name))
}
