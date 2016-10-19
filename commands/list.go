package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "List all credentials",
		Long: `The list command is used to list the saved credentials.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(List(args...))
		},
	})
}

func List(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	var accountNames []string
	for name := range safe.List() {
		accountNames = append(accountNames, name)
	}

	if len(accountNames) > 0 {
		sort.Strings(accountNames)
		for _, name := range accountNames {
			fmt.Println(name)
		}
	} else {
		fmt.Println("No accounts found")
	}

	return 0
}
