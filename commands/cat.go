package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "cat [name]",
		Short: "Cat a credential",
		Long:  "The cat command is used to print a credential to stdout.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Cat, cmd, args)
		},
	})
}

func Cat(args []string, flags *pflag.FlagSet) error {
	if len(args) != 1 {
		return &errors.InvalidCommandUsage{}
	}
	name := args[0]

	safe, err := loadSafe()
	if err != nil {
		return err
	}

	account, err := safe.Get(name)
	if err != nil {
		return err
	}

	fmt.Printf(`account:  %s
username: %s
password: %s
created:  %s
modified: %s
`, name, account.Username, account.Password,
		utils.FormatUnixTime(account.CreatedOn),
		utils.FormatUnixTime(account.ModifiedOn))
	return nil
}
