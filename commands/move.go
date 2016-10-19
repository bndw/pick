package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "mv [name] [new-name]",
		Short: "Rename a credential",
		Long: `The move command is used to rename a credential.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Move(args...))
		},
	})
}

func Move(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	name, newName, errCode := parseMoveArgs(args)
	if errCode > 0 {
		return errCode
	}

	if err := safe.Move(name, newName); err != nil {
		return handleError(err)
	}

	fmt.Println("Credential renamed")
	return 0
}

func parseMoveArgs(args []string) (name, newName string, errCode int) {
	if len(args) != 2 {
		fmt.Println("mv [name] [new-name]")
		return "", "", 1
	}

	name, newName = args[0], args[1]

	errCode = 0
	return
}
