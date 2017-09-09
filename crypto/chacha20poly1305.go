package crypto

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/bndw/pick/crypto/pbkdf2"
	"github.com/bndw/pick/crypto/scrypt"
	"github.com/bndw/pick/errors"
)

type ChaCha20Poly1305Client struct {
	settings      ChaCha20Poly1305Settings
	keyDerivation KeyDerivation
}

type ChaCha20Poly1305Settings struct {
	KeyDerivation string         `json:"keyderivation,omitempty" toml:"keyderivation"`
	PBKDF2        *pbkdf2.PBKDF2 `json:"pbkdf2,omitempty" toml:"pbkdf2"`
	Scrypt        *scrypt.Scrypt `json:"scrypt,omitempty" toml:"scrypt"`
}

type ChaCha20Poly1305Store struct {
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

const (
	chaCha20Poly1305DefaultKeyDerivation = keyDerivationTypePBKDF2
)

func DefaultChaCha20Poly1305Settings() *ChaCha20Poly1305Settings {
	return &ChaCha20Poly1305Settings{
		KeyDerivation: chaCha20Poly1305DefaultKeyDerivation,
		PBKDF2:        pbkdf2.New(),
		Scrypt:        scrypt.New(),
	}
}

func NewChaCha20Poly1305Client(settings *ChaCha20Poly1305Settings) (*ChaCha20Poly1305Client, error) {
	var kdf KeyDerivation
	switch settings.KeyDerivation {
	default:
		if settings.KeyDerivation != "" {
			fmt.Println("Invalid keyDerivation, using default")
		}
		fallthrough
	case keyDerivationTypePBKDF2:
		// Remove other settings
		// TODO(leon): This is shitty.
		settings.Scrypt = nil
		kdf = settings.PBKDF2
	case keyDerivationTypeScrypt:
		// Remove other settings
		// TODO(leon): This is shitty.
		settings.PBKDF2 = nil
		kdf = settings.Scrypt
	}
	return &ChaCha20Poly1305Client{
		settings:      *settings,
		keyDerivation: kdf,
	}, nil
}

func (c *ChaCha20Poly1305Client) keyLen() int {
	return chacha20poly1305.KeySize
}

func (c *ChaCha20Poly1305Client) deriveKey(password []byte, keyLen int) ([]byte, []byte, error) {
	return c.keyDerivation.DeriveKey(password, keyLen)
}

func (c *ChaCha20Poly1305Client) deriveKeyWithSalt(password, salt []byte, keyLen int) ([]byte, error) {
	return c.keyDerivation.DeriveKeyWithSalt(password, salt, keyLen)
}

func (c *ChaCha20Poly1305Client) Decrypt(data []byte, password []byte) ([]byte, error) {
	var store ChaCha20Poly1305Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	key, err := c.deriveKeyWithSalt(password, store.Salt, c.keyLen())
	if err != nil {
		return nil, err
	}

	cpc, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := cpc.Open(nil, store.Nonce, store.Ciphertext, nil)
	if err != nil {
		return nil, errors.ErrSafeDecryptionFailed
	}

	return plaintext, err
}

func (c *ChaCha20Poly1305Client) Encrypt(plaintext []byte, password []byte) (data []byte, err error) {
	key, salt, err := c.deriveKey(password, c.keyLen())
	if err != nil {
		return
	}

	cpc, err := chacha20poly1305.New(key)
	if err != nil {
		return
	}

	nonce := make([]byte, cpc.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}

	ciphertext := cpc.Seal(nil, nonce, plaintext, nil)

	store := ChaCha20Poly1305Store{
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}
	data, err = json.Marshal(store)
	if err != nil {
		return
	}

	return
}
