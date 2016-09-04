package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"math/big"
	pseudoRand "math/rand"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func init() {
	pseudoRand.Seed(time.Now().UTC().UnixNano())
}

// decrypt uses PGP to decrypt symmetrically encrypted and armored text
// with the provided password.
func Decrypt(text []byte, password []byte) (decryptedText []byte, err error) {
	decbuf := bytes.NewBuffer(text)

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

	decryptedText = decryptedBuf

	return
}

// encrypt uses PGP to symmetrically encrypt and armor text with the
// provided password.
func Encrypt(text []byte, password []byte) (encryptedText []byte, err error) {
	encbuf := bytes.NewBuffer(nil)

	w, err := armor.Encode(encbuf, "PGP SIGNATURE", nil)
	if err != nil {
		return
	}

	plaintext, err := openpgp.SymmetricallyEncrypt(w, password, nil, nil)
	if err != nil {
		return
	}

	_, err = plaintext.Write(text)

	plaintext.Close()
	w.Close()

	encryptedText = encbuf.Bytes()

	return
}

// GeneratePassword generates a password.
func GeneratePassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buffer := make([]byte, length)
	max := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		buffer[i] = chars[randomInt(max)]
	}

	return string(buffer)
}

func randomInt(max *big.Int) int {
	if rand, err := rand.Int(rand.Reader, max); err != nil {
		// Fallback to pseudo-random
		return pseudoRand.Intn(int(max.Int64()))
	} else {
		return int(rand.Int64())
	}
}
