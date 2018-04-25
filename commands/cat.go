package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/safe"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	cmd := &cobra.Command{
		Use:   "cat [name]",
		Short: "Cat a credential",
		Long:  "The cat command is used to print a credential to stdout.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Cat, cmd, args)
		},
	}
	cmd.Flags().Bool("history", false, "show credential history")
	rootCmd.AddCommand(cmd)
}

func printAccount(account *safe.Account, showHistory bool) {
	// Print header
	if showHistory && account.History != nil {
		history := account.History
		account.History = nil
		history = append(history, *account)
		for i, l := 0, len(history); i < l; i++ {
			printCredential(&history[i], "  ", i == 0)
		}
	} else {
		// Print a credential
		printCredential(account, "", true)
	}
}

func printCredential(account *safe.Account, printPrefix string, isInitialAccount bool) {
	// Print a credential
	var createdOrModified string
	if isInitialAccount {
		createdOrModified = "created"
	} else {
		createdOrModified = "modified"
	}
	fmt.Printf("%s: %s\n%susername: %s\n%spassword: %s\n",
		createdOrModified,
		utils.FormatUnixTime(account.ModifiedOn),
		printPrefix, account.Username,
		printPrefix, account.Password)
}

func Cat(args []string, flags *pflag.FlagSet) error {
	if len(args) != 1 {
		return errors.ErrInvalidCommandUsage
	}
	name := args[0]

	showHistory, err := parseCatFlags(flags)
	if err != nil {
		return err
	}

	safe, err := newSafeLoader(false).Load()
	if err != nil {
		return err
	}

	account, err := safe.Get(name)
	if err != nil {
		return err
	}

	printAccount(account, showHistory)
	return nil
}

func parseCatFlags(flags *pflag.FlagSet) (showHistory bool, err error) {
	showHistory, err = flags.GetBool("history")
	return
}
