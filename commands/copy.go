package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/strings"
	"github.com/bndw/pick/utils/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "cp [name]",
		Short: "Copy a credential to the clipboard",
		Long:  "The copy command is used to copy a credential's password to the clipboard.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Copy, cmd, args)
		},
	})
}

func Copy(args []string, flags *pflag.FlagSet) error {
	if len(args) != 1 {
		return errors.ErrInvalidCommandUsage
	}
	name := args[0]

	safe, err := newSafeLoader(false).Load()
	if err != nil {
		return err
	}

	account, err := safe.Get(name)
	if err != nil {
		return err
	}

	if err := clipboard.Copy(account.Password, safe.Config.General.Clipboard.ClearAfter); err != nil {
		return err
	}
	fmt.Println(strings.PasswordCopiedToClipboard)

	return nil
}
