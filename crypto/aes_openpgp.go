package crypto

import (
	"bytes"
	"crypto"
	"errors"
	"io/ioutil"
	"strings"

	_ "crypto/sha256"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type AESOpenPGPClient struct {
	packetConfig *packet.Config
}

func NewAESOpenPGPClient(config Config) (*AESOpenPGPClient, error) {

	// TODO(): Construct a packet from the config
	// pc := &packet.Config{
	// 	 DefaultHash: getHashFromConfig(config),
	// }

	return &AESOpenPGPClient{}, nil
}

// decrypt uses PGP to decrypt symmetrically encrypted and armored text
// with the provided password.
func (*AESOpenPGPClient) Decrypt(ciphertext []byte, password []byte) (plaintext []byte, err error) {
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
			return nil, errors.New("Unable to unlock safe with provided password")
		}

		failed = true

		return password, nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, nil, prompt, nil)

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
func (*AESOpenPGPClient) Encrypt(plaintext []byte, password []byte) (ciphertext []byte, err error) {
	encbuf := bytes.NewBuffer(nil)

	w, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return
	}
	defer w.Close()

	pt, err := openpgp.SymmetricallyEncrypt(w, password, nil, nil)
	if err != nil {
		return
	}
	defer pt.Close()

	if _, err := pt.Write(plaintext); err != nil {
		return nil, err
	}

	ciphertext = encbuf.Bytes()
	return
}

func getHashFromConfig(config Config) crypto.Hash {
	hash, ok := config.Settings["hash"].(string)
	if !ok {
		// No hash set, let the default case pick it up
	}

	switch strings.ToLower(hash) {
	default:
		return crypto.SHA256
	case "sha256":
		return crypto.SHA256
	}
}
