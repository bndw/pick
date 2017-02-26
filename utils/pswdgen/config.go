package pswdgen

import (
	"strings"
)

const (
	PasswordModeNonInteractive = "interactive"
	PasswordModeInteractive    = "nonInteractive"

	DefaultPasswordLength   = 25
	defaultPasswordStrength = 2 // alphanum
	defaultPasswordMode     = PasswordModeNonInteractive
)

type Config struct {
	Length   int `toml:"length"`
	Strength int
	Mode     string `toml:"mode"`
}

func NewDefaultConfig() Config {
	return Config{
		Length:   DefaultPasswordLength,
		Strength: defaultPasswordStrength,
		Mode:     defaultPasswordMode,
	}
}

func StrengthByString(s string) int {
	s = strings.ToLower(s)
	switch s {
	default:
		fallthrough
	case "full":
		return 3
	case "alphanum":
		return 2
	case "alpha":
		return 1
	case "num":
		return 0
	}
}
