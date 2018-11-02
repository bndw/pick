package commands

import (
	"errors"

	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	// TODO: This command is deprecated and will be removed soon.
	rootCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export decrypted credentials in JSON format",
		Long:  "The export command is used to export decrypted credentials in JSON format.",
		Run: func(cmd *cobra.Command, args []string) {
			runMovedCommand(Export, cmd, args, "safe export")
		},
		Hidden: true,
	})
}

func Export(args []string, flags *pflag.FlagSet) error {
	if !utils.Confirm("Do you really want to dump your whole pick safe?", false) {
		return errors.New("Aborted as requested")
	}

	safe, err := newSafeLoader(false).Load()
	if err != nil {
		return err
	}

	accounts := safe.List()
	if len(accounts) < 1 {
		return errors.New("No accounts to export")
	}

	utils.PrettyPrint(accounts)
	return nil
}
