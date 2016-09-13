package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "export",
		Short: "Export decrypted credentials in JSON format",
		Long: `The export command is used to export decrypted credentials in JSON
format.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Export(args...))
		},
	})

}

func Export(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	accounts := safe.List()
	if len(accounts) < 1 {
		fmt.Println("No accounts to export")
		return 1
	}

	utils.PrettyPrint(accounts)
	return 0
}
