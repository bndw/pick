package pswdgen

const (
	charsNum      = "0123456789"
	charsAlpha    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsAlphaNum = charsNum + charsAlpha
	charsFull     = charsAlphaNum + "-_.:,;()[]{}?!\"§$%&/=´`°^@|#'"

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
