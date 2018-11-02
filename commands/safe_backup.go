package commands

import (
	"fmt"

	"github.com/bndw/pick/safe"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	// TODO: This command is deprecated and will be removed soon.
	rootCmd.AddCommand(&cobra.Command{
		Use:   "backup",
		Short: "Backup the safe",
		Long:  "The backup command is used to backup your current safe.",
		Run: func(cmd *cobra.Command, args []string) {
			runMovedCommand(Backup, cmd, args, "safe backup")
		},
		Hidden: true,
	})
}

func Backup(args []string, flags *pflag.FlagSet) error {
	backendClient, err := newBackendClient()
	if err != nil {
		return err
	}

	if err := safe.Backup(backendClient); err != nil {
		return err
	}

	fmt.Println("Backup created")
	return nil
}
