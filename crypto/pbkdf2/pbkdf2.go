package pbkdf2

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2 struct {
	Hash       string `json:"hash,omitempty" toml:"hash"`
	Iterations int    `json:"iterations,omitempty" toml:"iterations"`
	SaltLen    int    `json:"saltlen,omitempty" toml:"saltlen"`
}

const (
	hashSHA256        = "sha256"
	hashSHA512        = "sha512"
	defaultHash       = hashSHA512
	defaultIterations = 100000
	defaultSaltLen    = 16
)

func New() *PBKDF2 {
	return &PBKDF2{
		Hash:       defaultHash,
		Iterations: defaultIterations,
		SaltLen:    defaultSaltLen,
	}
}

func (p *PBKDF2) HashFunc() func() hash.Hash {
	hash := p.Hash
	switch hash {
	default:
		if hash != "" {
			fmt.Println("Invalid PBKDF2 Hash, using default")
		}
		fallthrough
	case hashSHA512:
		return sha512.New
	case hashSHA256:
		return sha256.New
	}
}

func (p *PBKDF2) DeriveKeyWithSalt(password []byte, salt []byte, keyLen int) ([]byte, error) {
	iterations := p.Iterations
	hashFunc := p.HashFunc()
	return pbkdf2.Key(password, salt, iterations, keyLen, hashFunc), nil
}

func (p *PBKDF2) DeriveKey(password []byte, keyLen int) ([]byte, []byte, error) {
	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	key, _ := p.DeriveKeyWithSalt(password, salt, keyLen)

	return key, salt, nil
}
