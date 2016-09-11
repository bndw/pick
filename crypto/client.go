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
		return NewOpenPGPClient(*config.OpenPGPSettings)
	}
}
