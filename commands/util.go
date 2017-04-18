package commands

import (
	"bytes"
	builtinerrors "errors"
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
		if err == errors.ErrInvalidCommandUsage {
			cmd.Usage()
			os.Exit(1)
		}
		os.Exit(handleError(err))
	}
	os.Exit(0)
}

type safeLoader struct {
	password     *[]byte
	maxLoadTries int
	loadTries    int
}

func newSafeLoader() *safeLoader {
	return &safeLoader{
		maxLoadTries: 0,
	}
}

func (sl *safeLoader) RememberPassword() {
	sl.maxLoadTries++
}

func (sl *safeLoader) Load() (*safe.Safe, error) {
	backendClient, err := newBackendClient()
	if _, err := backendClient.Load(); err != nil {
		return nil, builtinerrors.New("pick not yet initialized. Please run the init command first")
	}
	if err != nil {
		return nil, err
	}
	return sl.LoadWithBackendClient(backendClient)
}

func (sl *safeLoader) LoadWithBackendClient(backendClient backends.Client) (*safe.Safe, error) {
	if sl.password == nil {
		password, err := utils.GetPasswordInput(fmt.Sprintf("Enter your master password for safe '%s'", backendClient.SafeLocation()))
		if err != nil {
			return nil, err
		}
		sl.password = &password
	}

	cryptoClient, err := newCryptoClient()
	if err != nil {
		return nil, err
	}

	s, err := safe.Load(
		*sl.password,
		backendClient,
		cryptoClient,
		config,
	)
	if err != nil {
		if sl.maxLoadTries > sl.loadTries {
			// Reset stored password and load again — asking for a new password
			sl.password = nil
			sl.loadTries++
			return sl.LoadWithBackendClient(backendClient)
		}
		return nil, err
	}
	return s, nil
}

func initSafe() error {
	backendClient, err := newBackendClient()
	if err != nil {
		return err
	}

	if _, err := backendClient.Load(); err == nil {
		return builtinerrors.New("pick was already initialized")
	}

	password, err := utils.GetPasswordInput("Please set a master password. This is the only password you need to remember")
	if err != nil {
		return err
	}
	passwordConfirm, err := utils.GetPasswordInput("Please confirm your master password")
	if err != nil {
		return err
	}
	if !bytes.Equal(password, passwordConfirm) {
		return builtinerrors.New("Master passwords do not match")
	}

	cryptoClient, err := newCryptoClient()
	if err != nil {
		return err
	}

	s, err := safe.Load(
		password,
		backendClient,
		cryptoClient,
		config,
	)
	if err != nil {
		return err
	}

	if err := s.Init(); err != nil {
		return err
	}

	fmt.Println("pick initialized")
	return nil
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
