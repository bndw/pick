package crypto

type Config struct {
	Type             string `json:"type" toml:"type"`
	*OpenPGPSettings `json:"openpgp,omitempty" toml:"openpgp"`
	*AESGCMSettings  `json:"aes_gcm,omitempty" toml:"aes_gcm"`
}

const (
	ConfigTypeOpenPGP = "openpgp"
	ConfigTypeAESGCM  = "aes_gcm"
	cipherAES256      = "aes256"
	cipherAES128      = "aes128"
	hashSHA256        = "sha256"
	hashSHA512        = "sha512"
)

func NewDefaultConfig() Config {
	return Config{
		Type:            ConfigTypeOpenPGP,
		OpenPGPSettings: DefaultOpenPGPSettings(),
		AESGCMSettings:  DefaultAESGCMSettings(),
	}
}

func NewDefaultConfigWithType(t string) Config {
	dc := NewDefaultConfig()
	dc.Type = t
	return dc
}
