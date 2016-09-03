package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

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
func GeneratePassword(length int) (password string, err error) {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

	passwordBuf := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	charLen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0

	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			return "", err
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}

			passwordBuf[i] = chars[c%charLen]
			i++

			if i == length {
				// We're done
				return string(passwordBuf), nil
			}
		}
	}

	// noop
	return "", errors.New("Unable to generate password")
}
