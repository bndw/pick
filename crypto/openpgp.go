package crypto

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/bndw/pick/errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type OpenPGPClient struct {
	settings     OpenPGPSettings
	packetConfig *packet.Config
}

type OpenPGPSettings struct {
	Cipher   string `json:"cipher,omitempty" toml:"cipher"`
	S2KCount int    `json:"s2kcount,omitempty" toml:"s2kcount"`
}

const (
	openpgpDefaultCipher   = cipherAES256
	openpgpDefaultS2KCount = 65011712
)

func DefaultOpenPGPSettings() *OpenPGPSettings {
	return &OpenPGPSettings{
		Cipher:   openpgpDefaultCipher,
		S2KCount: openpgpDefaultS2KCount,
	}
}

func NewOpenPGPClient(settings *OpenPGPSettings) (*OpenPGPClient, error) {
	c := &OpenPGPClient{
		settings: *settings,
	}
	c.packetConfig = &packet.Config{
		DefaultCipher: c.cipherFunc(),
		S2KCount:      c.s2kCount(),
	}
	return c, nil
}

func (c *OpenPGPClient) cipherFunc() packet.CipherFunction {
	switch c.settings.Cipher {
	default:
		if c.settings.Cipher != "" {
			fmt.Println("Invalid cipher, using default")
		}
		fallthrough
	case cipherAES256:
		return packet.CipherAES256
	case cipherAES128:
		return packet.CipherAES128
	}
}

func (c *OpenPGPClient) s2kCount() int {
	return c.settings.S2KCount
}

// decrypt uses PGP to decrypt symmetrically encrypted and armored text
// with the provided password.
func (c *OpenPGPClient) Decrypt(ciphertext []byte, password []byte) (plaintext []byte, err error) {
	decbuf := bytes.NewBuffer(ciphertext)

	armorBlock, err := armor.Decode(decbuf)
	if err != nil {
		return
	}

	failed := false
	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		// If the given passphrase isn't correct, the function will be called again, forever.
		// This method will fail fast.
		// Ref: https://godoc.org/golang.org/x/crypto/openpgp#PromptFunction
		if failed {
			return nil, errors.ErrSafeDecryptionFailed
		}

		failed = true

		return password, nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, nil, prompt, c.packetConfig)

	if err != nil {
		return
	}

	decryptedBuf, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return
	}

	plaintext = decryptedBuf
	return
}

// encrypt uses PGP to symmetrically encrypt and armor text with the
// provided password.
func (c *OpenPGPClient) Encrypt(plaintext []byte, password []byte) (ciphertext []byte, err error) {
	encbuf := bytes.NewBuffer(nil)

	w, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return
	}
	defer w.Close() // nolint: errcheck

	pt, err := openpgp.SymmetricallyEncrypt(w, password, nil, c.packetConfig)
	if err != nil {
		return
	}
	defer pt.Close() // nolint: errcheck

	if _, err := pt.Write(plaintext); err != nil {
		return nil, err
	}

	// Force-close writer to flush their cache
	if err = pt.Close(); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}

	ciphertext = encbuf.Bytes()
	return
}
