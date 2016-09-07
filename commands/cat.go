package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "cat [name]",
		Short: "Cat a credential",
		Long: `The cat command is used to print a credential to stdout.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("USAGE: cat [name]")
				os.Exit(1)
			}

			os.Exit(Cat(args...))
		},
	})
}

func Cat(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	account, err := safe.Get(args[0])
	if err != nil {
		return handleError(err)
	}

	fmt.Printf(`account:  %s
username: %s
password: %s
created:  %s
`, account.Name, account.Username, account.Password,
		utils.FormatUnixTime(account.CreatedOn))
	return 0
}
