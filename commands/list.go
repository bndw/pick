package commands

import (
	"errors"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "ls",
		Short: "List all accounts",
		Long:  "The list command is used to list the saved accounts.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(List, cmd, args)
		},
	})
}

func List(args []string, flags *pflag.FlagSet) error {
	safe, err := newSafeLoader(false).Load()
	if err != nil {
		return err
	}

	var accountNames []string
	for name := range safe.List() {
		accountNames = append(accountNames, name)
	}

	if len(accountNames) == 0 {
		return errors.New("No accounts found")
	}

	sort.Strings(accountNames)
	for _, name := range accountNames {
		fmt.Println(name)
	}
	return nil
}
