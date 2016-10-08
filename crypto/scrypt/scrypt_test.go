package scrypt_test

import (
	"testing"

	"github.com/bndw/pick/crypto/scrypt"
)

const (
	pass    = "my-secret"
	salt    = "abcd"
	keyLen  = 2
	saltLen = 4
	n       = 2
	r       = 1
	p       = 1
)

var s = NewWeakScrypt()

func NewWeakScrypt() *scrypt.Scrypt {
	s := scrypt.New()
	s.SaltLen = saltLen
	s.N = n
	s.R = r
	s.P = p
	return s
}

func TestDeriveKeyWithSalt(t *testing.T) {
	key, err := s.DeriveKeyWithSalt([]byte(pass), []byte(salt), keyLen)
	if err != nil {
		t.Error(err)
	}
	if len(key) != keyLen {
		t.Errorf("Unexpected key len: got %d, want %d", len(key), keyLen)
	}
}

func TestDeriveKey(t *testing.T) {
	key, salt, err := s.DeriveKey([]byte(pass), keyLen)
	if err != nil {
		t.Error(err)
	}
	if len(key) != keyLen {
		t.Errorf("Unexpected key len: got %d, want %d", len(key), keyLen)
	}
	if len(salt) != saltLen {
		t.Errorf("Unexpected salt len: got %d, want %d", len(salt), saltLen)
	}
}
