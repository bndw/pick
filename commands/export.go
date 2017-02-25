package commands

import (
	"errors"

	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export decrypted credentials in JSON format",
		Long:  "The export command is used to export decrypted credentials in JSON format.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Export, cmd, args)
		},
	})

}

func Export(args []string, flags *pflag.FlagSet) error {
	safe, err := newSafeLoader().Load()
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
