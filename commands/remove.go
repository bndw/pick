package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "rm [name]",
		Short: "Remove a credential",
		Long:  "The remove command is used to remove a saved credential.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Remove, cmd, args)
		},
	})
}

func Remove(args []string, flags *pflag.FlagSet) error {
	if len(args) != 1 {
		return errors.ErrInvalidCommandUsage
	}
	name := args[0]

	safe, err := newSafeLoader().Load()
	if err != nil {
		return err
	}

	if err := safe.Remove(name); err != nil {
		return err
	}

	fmt.Println("Credential removed")
	return nil
}
