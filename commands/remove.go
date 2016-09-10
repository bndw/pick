package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "rm [name]",
		Short: "Remove a credential",
		Long: `The remove command is used to remove a saved credential.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Remove(args...))
		},
	})
}

func Remove(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	if err := safe.Remove(args[0]); err != nil {
		return handleError(err)
	}

	fmt.Println("Credential removed")
	return 0
}
