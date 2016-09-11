package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"hash"

	"github.com/bndw/pick/errors"
	"golang.org/x/crypto/pbkdf2"
)

type AESGCMClient struct {
	settings AESGCMSettings
}

type AESGCMSettings struct {
	KeyLen           int    `json:"keylen,omitempty" toml:"keylen"`
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
	aesGCMDefaultKeyLen           = 32
	aesGCMDefaultPbkdf2Hash       = hashSHA512
	aesGCMDefaultPbkdf2Iterations = 100000
	aesGCMDefaultPbkdf2SaltLen    = 16
)

func DefaultAESGCMSettings() *AESGCMSettings {
	return &AESGCMSettings{
		KeyLen:           aesGCMDefaultKeyLen,
		Pbkdf2Hash:       aesGCMDefaultPbkdf2Hash,
		Pbkdf2Iterations: aesGCMDefaultPbkdf2Iterations,
		Pbkdf2SaltLen:    aesGCMDefaultPbkdf2SaltLen,
	}
}

func NewAESGCMClient(settings AESGCMSettings) (*AESGCMClient, error) {
	return &AESGCMClient{
		settings: settings,
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
	case 16:
	case 24:
	case 32:
	}
	return keyLen
}

func (c *AESGCMClient) pbkdf2HashFunc() func() hash.Hash {
	pbkdf2Hash := c.settings.Pbkdf2Hash
	switch pbkdf2Hash {
	default:
		if pbkdf2Hash != "" {
			fmt.Println("Invalid PBKDF2 Hash, using default")
		}
		fallthrough
	case hashSHA512:
		return sha512.New
	}
}

func (c *AESGCMClient) pbkdf2SaltLen() int {
	return c.settings.Pbkdf2SaltLen
}

func (c *AESGCMClient) pbkdf2Iterations() int {
	return c.settings.Pbkdf2Iterations
}

func (c *AESGCMClient) deriveKeyWithSalt(password []byte, salt []byte) []byte {
	keyLen := c.keyLen()
	pbkdf2Iterations := c.pbkdf2Iterations()
	pbkdf2HashFunc := c.pbkdf2HashFunc()
	return pbkdf2.Key(password, salt, pbkdf2Iterations, keyLen, pbkdf2HashFunc)
}

func (c *AESGCMClient) deriveKey(password []byte) ([]byte, []byte, error) {
	salt := make([]byte, c.pbkdf2SaltLen())
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	return c.deriveKeyWithSalt(password, salt), salt, nil
}

func (c *AESGCMClient) Decrypt(data []byte, password []byte) (plaintext []byte, err error) {
	var store AESGCMStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	key := c.deriveKeyWithSalt(password, store.Salt)

	ac, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(ac)
	if err != nil {
		return
	}

	plaintext, err = gcm.Open(nil, store.Nonce, store.Ciphertext, nil)
	if err != nil {
		return nil, &errors.SafeDecryptionFailed{}
	}

	return
}

func (c *AESGCMClient) Encrypt(plaintext []byte, password []byte) (data []byte, err error) {
	key, salt, err := c.deriveKey(password)
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
