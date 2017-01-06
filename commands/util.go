package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/safe"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func runCommand(c func([]string, *pflag.FlagSet) error, cmd *cobra.Command, args []string) {
	if err := c(args, cmd.Flags()); err != nil {
		if _, isUsageErr := err.(*errors.InvalidCommandUsage); isUsageErr {
			cmd.Usage()
			os.Exit(1)
		}
		os.Exit(handleError(err))
	}
	os.Exit(0)
}

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
		config,
	)
}

func newBackendClient() (backends.Client, error) {
	return backends.New(&config.Storage)
}

func newCryptoClient() (crypto.Client, error) {
	return crypto.New(&config.Encryption)
}

func handleError(err error) int {
	fmt.Println(err)
	return 1
}

func isInvalidCommandUsage(err error) bool {
	_, ok := err.(*errors.InvalidCommandUsage)
	return ok
}
