package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "mv [name] [new-name]",
		Short: "Rename an account",
		Long:  "The move command is used to rename an account.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Move, cmd, args)
		},
	})
}

func Move(args []string, flags *pflag.FlagSet) error {
	safe, err := newSafeLoader(true).Load()
	if err != nil {
		return err
	}

	name, newName, err := parseMoveArgs(args)
	if err != nil {
		return err
	}

	if err := safe.Move(name, newName); err != nil {
		return err
	}

	fmt.Println("Account renamed")
	return nil
}

func parseMoveArgs(args []string) (name, newName string, err error) {
	if len(args) != 2 {
		err = errors.ErrInvalidCommandUsage
		return
	}

	name, newName = args[0], args[1]

	return
}
