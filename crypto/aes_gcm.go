package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/bndw/pick/crypto/pbkdf2"
	"github.com/bndw/pick/crypto/scrypt"
	"github.com/bndw/pick/errors"
)

type AESGCMClient struct {
	settings      AESGCMSettings
	keyDerivation KeyDerivation
}

type AESGCMSettings struct {
	KeyLen        int            `json:"keylen,omitempty" toml:"keylen"`
	KeyDerivation string         `json:"keyderivation,omitempty" toml:"keyderivation"`
	PBKDF2        *pbkdf2.PBKDF2 `json:"pbkdf2,omitempty" toml:"pbkdf2"`
	Scrypt        *scrypt.Scrypt `json:"scrypt,omitempty" toml:"scrypt"`
	// Warning: Deprecated. These three Pbkdf2 configs are required for backwards-compatibility :(
	Pbkdf2Hash       string `json:"pbkdf2hash,omitempty" toml:"pbkdf2hash"`
	Pbkdf2Iterations int    `json:"pbkdf2iterations,omitempty" toml:"pbkdf2iterations"`
	Pbkdf2SaltLen    int    `json:"pbkdf2saltlen,omitempty" toml:"pbkdf2saltlen"`
}

type AESGCMStore struct {
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

const (
	aesGCMDefaultKeyLen        = cipherLenAES256
	aesGCMDefaultKeyDerivation = keyDerivationTypePBKDF2
)

func DefaultAESGCMSettings() *AESGCMSettings {
	return &AESGCMSettings{
		KeyLen:        aesGCMDefaultKeyLen,
		KeyDerivation: aesGCMDefaultKeyDerivation,
		PBKDF2:        pbkdf2.New(),
		Scrypt:        scrypt.New(),
	}
}

func NewAESGCMClient(settings *AESGCMSettings) (*AESGCMClient, error) {
	if settings.PBKDF2 == nil {
		// Probably a safe which uses the old config, backwards-compatibility mode
		settings.PBKDF2 = pbkdf2.New()
		settings.PBKDF2.Hash = settings.Pbkdf2Hash
		settings.PBKDF2.Iterations = settings.Pbkdf2Iterations
		settings.PBKDF2.SaltLen = settings.Pbkdf2SaltLen
	}
	var kdf KeyDerivation
	switch settings.KeyDerivation {
	default:
		if settings.KeyDerivation != "" {
			fmt.Println("Invalid keyDerivation, using default")
		}
		fallthrough
	case keyDerivationTypePBKDF2:
		// Remove other settings
		// TODO: Remove other settings in a more elegant way
		settings.Scrypt = nil
		kdf = settings.PBKDF2
	case keyDerivationTypeScrypt:
		// Remove other settings
		// TODO: Remove other settings in a more elegant way
		settings.PBKDF2 = nil
		kdf = settings.Scrypt
	}
	return &AESGCMClient{
		settings:      *settings,
		keyDerivation: kdf,
	}, nil
}

func (c *AESGCMClient) keyLen() int {
	keyLen := c.settings.KeyLen
	switch keyLen {
	default:
		if keyLen != 0 {
			fmt.Println("Invalid keyLen, using default")
		}
		return aesGCMDefaultKeyLen
	case cipherLenAES128:
	case cipherLenAES192:
	case cipherLenAES256:
	}
	return keyLen
}

func (c *AESGCMClient) deriveKey(password []byte, keyLen int) ([]byte, []byte, error) {
	return c.keyDerivation.DeriveKey(password, keyLen)
}

func (c *AESGCMClient) deriveKeyWithSalt(password, salt []byte, keyLen int) ([]byte, error) {
	return c.keyDerivation.DeriveKeyWithSalt(password, salt, keyLen)
}

func (c *AESGCMClient) Decrypt(data []byte, password []byte) ([]byte, error) {
	var store AESGCMStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	key, err := c.deriveKeyWithSalt(password, store.Salt, c.keyLen())
	if err != nil {
		return nil, err
	}

	ac, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(ac)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, store.Nonce, store.Ciphertext, nil)
	if err != nil {
		return nil, errors.ErrSafeDecryptionFailed
	}

	return plaintext, nil
}

func (c *AESGCMClient) Encrypt(plaintext []byte, password []byte) (data []byte, err error) {
	key, salt, err := c.deriveKey(password, c.keyLen())
	if err != nil {
		return
	}

	ac, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(ac)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	store := AESGCMStore{
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
