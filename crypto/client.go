package crypto

import (
	"fmt"
)

type Client interface {
	Decrypt(data, password []byte) (plaintext []byte, err error)
	Encrypt(plaintext, password []byte) (data []byte, err error)
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
		return NewOpenPGPClient(*config.OpenPGPSettings)
	case ConfigTypeAESGCM:
		// Remove other settings
		// TODO(leon): This is shitty.
		config.OpenPGPSettings = nil
		return NewAESGCMClient(*config.AESGCMSettings)
	}
}
