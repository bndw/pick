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

	backendClient, err := backends.New(config.Storage)
	if err != nil {
		return nil, err
	}

	cryptoClient, err := crypto.New(config.Encryption)
	if err != nil {
		return nil, err
	}

	return safe.Load(
		password,
		backendClient,
		cryptoClient,
	)
}

func handleError(err error) int {
	fmt.Println(err)
	return 1
}
