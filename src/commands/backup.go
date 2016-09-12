package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/safe"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "backup",
		Short: "Backup the safe",
		Long: `The backup command is used to backup your current safe.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Backup(args...))
		},
	})
}

func Backup(args ...string) int {
	backendClient, err := newBackendClient()
	if err != nil {
		return handleError(err)
	}

	if err := safe.Backup(backendClient); err != nil {
		return handleError(err)
	}

	fmt.Println("Backup created")
	return 0
}
