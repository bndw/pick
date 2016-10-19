package crypto

type Config struct {
	Type                      string `json:"type" toml:"type"`
	*OpenPGPSettings          `json:"openpgp,omitempty" toml:"openpgp"`
	*AESGCMSettings           `json:"aes_gcm,omitempty" toml:"aes_gcm"`
	*ChaCha20Poly1305Settings `json:"chachapoly,omitempty" toml:"chachapoly"`
}

const (
	ConfigTypeOpenPGP       = "openpgp"
	ConfigTypeAESGCM        = "aes_gcm"
	ConfigTypeChaChaPoly    = "chachapoly"
	keyDerivationTypePBKDF2 = "pbkdf2"
	keyDerivationTypeScrypt = "scrypt"
	cipherAES256            = "aes256"
	cipherAES128            = "aes128"
	cipherLenAES256         = 32
	cipherLenAES192         = 24
	cipherLenAES128         = 16
)

func NewDefaultConfig() Config {
	return Config{
		Type:                     ConfigTypeOpenPGP,
		OpenPGPSettings:          DefaultOpenPGPSettings(),
		AESGCMSettings:           DefaultAESGCMSettings(),
		ChaCha20Poly1305Settings: DefaultChaCha20Poly1305Settings(),
	}
}

func NewDefaultConfigWithType(t string) Config {
	dc := NewDefaultConfig()
	dc.Type = t
	return dc
}
