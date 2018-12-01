package crypto

import (
	"fmt"
)

type Client interface {
	Decrypt(data, password []byte) (plaintext []byte, err error)
	Encrypt(plaintext, password []byte) (data []byte, err error)
}

type KeyDerivation interface {
	DeriveKey(password []byte, keyLen int) (key []byte, salt []byte, err error)
	DeriveKeyWithSalt(password, salt []byte, keyLen int) (key []byte, err error)
}

func New(config *Config) (Client, error) {
	switch t := config.Type; t {
	default:
		config.Type = DefaultConfigType
		fmt.Printf("Invalid encryption type %q, using default %q\n", t, DefaultConfigType)
		// This won't recurse indefinitely as we have provided a valid config type above
		return New(config)
	case ConfigTypeChaChaPoly:
		// Remove other settings
		// TODO: Remove other settings in a more elegant way
		config.OpenPGPSettings = nil
		config.AESGCMSettings = nil
		return NewChaCha20Poly1305Client(config.ChaCha20Poly1305Settings)
	case ConfigTypeOpenPGP:
		// Remove other settings
		// TODO: Remove other settings in a more elegant way
		config.AESGCMSettings = nil
		config.ChaCha20Poly1305Settings = nil
		return NewOpenPGPClient(config.OpenPGPSettings)
	case ConfigTypeAESGCM:
		// Remove other settings
		// TODO: Remove other settings in a more elegant way
		config.OpenPGPSettings = nil
		config.ChaCha20Poly1305Settings = nil
		return NewAESGCMClient(config.AESGCMSettings)
	}
}
