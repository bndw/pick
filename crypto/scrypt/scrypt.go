package scrypt

import (
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

type Scrypt struct {
	SaltLen int `json:"saltlen,omitempty" toml:"saltlen"`
	N       int `json:"n,omitempty" toml:"n"`
	R       int `json:"r,omitempty" toml:"r"`
	P       int `json:"p,omitempty" toml:"p"`
}

const (
	defaultSaltLen = 16
	defaultN       = 16384
	defaultR       = 8
	defaultP       = 1
)

func New() *Scrypt {
	return &Scrypt{
		SaltLen: defaultSaltLen,
		N:       defaultN,
		R:       defaultR,
		P:       defaultP,
	}
}

func (s *Scrypt) DeriveKeyWithSalt(password, salt []byte, keyLen int) ([]byte, error) {
	key, err := scrypt.Key(password, salt, s.N, s.R, s.P, keyLen)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *Scrypt) DeriveKey(password []byte, keyLen int) ([]byte, []byte, error) {
	salt := make([]byte, s.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	key, err := s.DeriveKeyWithSalt(password, salt, keyLen)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
