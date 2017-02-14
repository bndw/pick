package pswdgen

const (
	defaultPasswordLength   = 25
	defaultPasswordStrength = 2 // alphanum
	defaultPasswordMode     = "nonInteractive"

	passwordModeInteractive = "interactive"
)

type Config struct {
	Length   int `toml:"length"`
	Strength int
	Mode     string `toml:"mode"`
}

func NewDefaultConfig() Config {
	return Config{
		Length:   defaultPasswordLength,
		Strength: defaultPasswordStrength,
		Mode:     defaultPasswordMode,
	}
}
