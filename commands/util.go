package commands

import (
	"fmt"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/safe"
	"github.com/bndw/pick/utils"
)

func loadSafe() (*safe.Safe, error) {
	password, err := utils.GetPasswordInput("Enter your master password")
	if err != nil {
		return nil, err
	}

	backendClient, err := newBackendClient()
	if err != nil {
		return nil, err
	}

	cryptoClient, err := newCryptoClient()
	if err != nil {
		return nil, err
	}

	return safe.Load(
		password,
		backendClient,
		cryptoClient,
	)
}

func newBackendClient() (backends.Client, error) {
	return backends.New(Config.Storage)
}

func newCryptoClient() (crypto.Client, error) {
	return crypto.New(Config.Encryption)
}

func handleError(err error) int {
	fmt.Println(err)
	return 1
}
