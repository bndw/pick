package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "cp [name]",
		Short: "Copy a credential to the clipboard",
		Long: `The copy command is used to copy a credential's password
to the clipboard.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("USAGE: copy [name]")
				os.Exit(1)
			}

			os.Exit(Copy(args...))
		},
	})
}

func Copy(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	if err := safe.Copy(args[0]); err != nil {
		return handleError(err)
	}

	return 0
}
