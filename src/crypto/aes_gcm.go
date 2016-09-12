package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/bndw/pick/crypto/pbkdf2"
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
	// These three Pbkdf2 configs are required for backwards-compatiblity :(
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
	aesGCMDefaultKeyLen        = 32
	aesGCMDefaultKeyDerivation = "pbkdf2"
)

func DefaultAESGCMSettings() *AESGCMSettings {
	return &AESGCMSettings{
		KeyLen:        aesGCMDefaultKeyLen,
		KeyDerivation: aesGCMDefaultKeyDerivation,
		PBKDF2:        pbkdf2.New(),
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
	case "pbkdf2":
		kdf = settings.PBKDF2
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
	case 16:
	case 24:
	case 32:
	}
	return keyLen
}

func (c *AESGCMClient) deriveKey(password []byte, keyLen int) ([]byte, []byte, error) {
	return c.keyDerivation.DeriveKey(password, keyLen)
}

func (c *AESGCMClient) deriveKeyWithSalt(password, salt []byte, keyLen int) ([]byte, error) {
	return c.keyDerivation.DeriveKeyWithSalt(password, salt, keyLen)
}

func (c *AESGCMClient) Decrypt(data []byte, password []byte) (plaintext []byte, err error) {
	var store AESGCMStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	key, err := c.deriveKeyWithSalt(password, store.Salt, c.keyLen())
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

	plaintext, err = gcm.Open(nil, store.Nonce, store.Ciphertext, nil)
	if err != nil {
		return nil, &errors.SafeDecryptionFailed{}
	}

	return
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
