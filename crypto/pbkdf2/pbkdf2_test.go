package pbkdf2_test

import (
	"testing"

	"github.com/bndw/pick/crypto/pbkdf2"
)

const (
	pass       = "my-secret"
	salt       = "abcd"
	keyLen     = 2
	hash       = "sha256"
	iterations = 10
	saltLen    = 4
)

var p = NewWeakPBKDF2()

func NewWeakPBKDF2() *pbkdf2.PBKDF2 {
	p := pbkdf2.New()
	p.Hash = hash
	p.Iterations = iterations
	p.SaltLen = saltLen
	return p
}

func TestDeriveKeyWithSalt(t *testing.T) {
	key, err := p.DeriveKeyWithSalt([]byte(pass), []byte(salt), keyLen)
	if err != nil {
		t.Error(err)
	}
	if len(key) != keyLen {
		t.Errorf("Unexpected key len: got %d, want %d", len(key), keyLen)
	}
}

func TestDeriveKey(t *testing.T) {
	key, salt, err := p.DeriveKey([]byte(pass), keyLen)
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
