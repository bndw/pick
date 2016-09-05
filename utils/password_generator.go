package utils

import (
	"crypto/rand"
	"math/big"
)

// GeneratePassword generates a password.
func GeneratePassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buffer := make([]byte, length)
	max := big.NewInt(int64(len(chars)))

	var index int
	var err error
	for i := 0; i < length; i++ {
		index, err = randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = chars[index]
	}

	return string(buffer), nil
}

func randomInt(max *big.Int) (int, error) {
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(rand.Int64()), nil
}
