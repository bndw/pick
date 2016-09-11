package crypto

type Config struct {
	Type             string `json:"type" toml:"type"`
	*OpenPGPSettings `json:"openpgp,omitempty" toml:"openpgp"`
}

const (
	ConfigTypeOpenPGP = "openpgp"
	cipherAES256      = "aes256"
	cipherAES128      = "aes128"
)

func NewDefaultConfig() Config {
	return Config{
		Type:            ConfigTypeOpenPGP,
		OpenPGPSettings: DefaultOpenPGPSettings(),
	}
}

func NewDefaultConfigWithType(t string) Config {
	dc := NewDefaultConfig()
	dc.Type = t
	return dc
}
