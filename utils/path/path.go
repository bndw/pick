package path

import (
	"strings"
)

const (
	modCurrent  = "."
	modParent   = ".."
	modSlashNix = "/"  // *nix
	modSlashWin = "\\" // Windows
)

// IsAbs returns true if the path provided is absolute, false otherwise.
func IsAbs(p string) bool {
	fs := []string{
		modSlashNix,
		modSlashWin,
	}
	for _, f := range fs {
		if strings.HasPrefix(p, f) {
			return true
		}
	}
	return false
}

// IsRel returns true is the path starts with either ".*" or "..*"
// where "*" is the OS-specific path delimiter, false otherwise.
func IsRel(p string) bool {
	fs := []string{
		modCurrent + modSlashNix,
		modParent + modSlashNix,

		modCurrent + modSlashWin,
		modParent + modSlashWin,
	}
	for _, f := range fs {
		if strings.HasPrefix(p, f) {
			return true
		}
	}
	return false
}

func TrimModPrefix(p string) string {
	mods := []string{
		modCurrent + modSlashNix,
		modParent + modSlashNix,

		modCurrent + modSlashWin,
		modParent + modSlashWin,

		modSlashNix,
		modSlashWin,
	}

	for {
		tmp := p
		for _, m := range mods {
			tmp = strings.TrimPrefix(tmp, m)
		}
		if tmp == p {
			break
		}
		p = tmp
	}
	return p
}
