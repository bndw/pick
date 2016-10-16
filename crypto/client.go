package crypto

import (
	"fmt"
)

type Client interface {
	Decrypt(data, password []byte) (plaintext []byte, err error)
	Encrypt(plaintext, password []byte) (data []byte, err error)
}

type KeyDerivation interface {
	DeriveKey(password []byte, keyLen int) ([]byte, []byte, error)
	DeriveKeyWithSalt(password, salt []byte, keyLen int) ([]byte, error)
}

func New(config *Config) (Client, error) {
	switch config.Type {
	default:
		if config.Type != "" {
			fmt.Println("Invalid encryption type, using default")
		}
		fallthrough
	case ConfigTypeOpenPGP:
		// Remove other settings
		// TODO(leon): This is shitty.
		config.AESGCMSettings = nil
		config.ChaCha20Poly1305Settings = nil
		return NewOpenPGPClient(config.OpenPGPSettings)
	case ConfigTypeAESGCM:
		// Remove other settings
		// TODO(leon): This is shitty.
		config.OpenPGPSettings = nil
		config.ChaCha20Poly1305Settings = nil
		return NewAESGCMClient(config.AESGCMSettings)
	case ConfigTypeChaChaPoly:
		// Remove other settings
		// TODO(leon): This is shitty.
		config.OpenPGPSettings = nil
		config.AESGCMSettings = nil
		return NewChaCha20Poly1305Client(config.ChaCha20Poly1305Settings)
	}
}
