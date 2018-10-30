package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	safeCmd := &cobra.Command{
		Use:   "safe --help",
		Short: "Perform operations on safe",
		Long:  "The safe command is used to perform operations on your pick safe.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	safeCmd.AddCommand(&cobra.Command{
		Use:   "backup",
		Short: "Backup the safe",
		Long:  "The backup command is used to backup your current safe.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Backup, cmd, args)
		},
	})
	safeCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export decrypted credentials in JSON format",
		Long:  "The export command is used to export decrypted credentials in JSON format.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Export, cmd, args)
		},
	})
	safeCmd.AddCommand(&cobra.Command{
		Use:   "sync [other-pick-safe]",
		Short: "Sync current safe with another pick safe",
		Long:  "The sync command is used to sync the current pick safe with another pick safe.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Sync, cmd, args)
		},
	})

	rootCmd.AddCommand(safeCmd)
}
