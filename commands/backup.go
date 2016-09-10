package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
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

	if err := backendClient.Backup(); err != nil {
		return handleError(err)
	}

	fmt.Println("Backup created")
	return 0
}
