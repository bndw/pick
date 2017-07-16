package randomstring

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	// CharsNum contains numbers from 0-9
	CharsNum = "0123456789"
	// CharsAlpha contains the full English alphabet: letters a-z and A-Z
	CharsAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// CharsAlphaNum is a combination of CharsNum and CharsAlpha
	CharsAlphaNum = CharsNum + CharsAlpha
	// CharsASCII contains all printable ASCII characters in code range [32, 126]
	CharsASCII = CharsAlphaNum + " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var (
	// ErrInvalidLengthSpecified is returned when the length specified is invalid
	ErrInvalidLengthSpecified = errors.New("invalid password length specified")
)

// Generate generates a cryptographically secure and unbiased string of length 'l' using alphabet 'dict'
func Generate(l int, dict string) (string, error) {
	// Length needs to be in range [1, 1<<31-1]
	if l <= 0 || l > 1<<31-1 {
		return "", ErrInvalidLengthSpecified
	}

	buf := make([]byte, l)
	max := big.NewInt(int64(len(dict)))

	for i := 0; i < l; i++ {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}

		buf[i] = dict[index]
	}

	return string(buf), nil
}

func randomInt(max *big.Int) (int, error) {
	i, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(i.Int64()), nil
}
