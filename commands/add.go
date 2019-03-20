package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/strings"
	"github.com/bndw/pick/utils"
	"github.com/bndw/pick/utils/clipboard"
	"github.com/bndw/pick/utils/pswdgen"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "add [name] [username]",
		Short: "Add an account",
		Long:  "The add command is used to add a new account.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Add, cmd, args)
		},
	})
}

func Add(args []string, flags *pflag.FlagSet) error {
	safe, err := newSafeLoader(true).Load()
	if err != nil {
		return err
	}

	name, username, err := parseAddArgs(args)
	if err != nil {
		return err
	}

	var password string
	if utils.Confirm("Generate password", true) {
		if password, err = pswdgen.Generate(config.General.Password); err != nil {
			return err
		}
	} else {
		var tmp []byte
		if tmp, err = utils.GetPasswordInput(fmt.Sprintf("Enter a password for %s", name)); err != nil {
			return err
		}
		password = string(tmp)
	}

	account, err := safe.Add(name, username, password)
	if err == errors.ErrAccountAlreadyExists && overwrite(name) {
		var editErr error
		if account, editErr = safe.Edit(name, username, password); editErr != nil {
			return editErr
		}
	} else if err != nil {
		return err
	}

	fmt.Println("Account added")
	if utils.Confirm("Copy password to clipboard", true) {
		if err := clipboard.Copy(account.Password, safe.Config.General.Clipboard.ClearAfter); err != nil {
			return err
		}
		fmt.Println(strings.PasswordCopiedToClipboard)
	}
	return nil
}

func overwrite(name string) bool {
	prompt := fmt.Sprintf("%s already exists, overwrite", name)
	return utils.Confirm(prompt, false)
}

func parseAddArgs(args []string) (name, username string, err error) {
	if len(args) > 2 {
		err = errors.ErrInvalidCommandUsage
		return
	}

	switch len(args) {
	case 2:
		username = args[1]
		fallthrough
	case 1:
		name = args[0]
	}

	if name == "" {
		if name, err = utils.GetInput("Enter an account name"); err != nil {
			return
		}
	}

	if username == "" {
		if username, err = utils.GetInput(fmt.Sprintf("Enter a username for %s", name)); err != nil {
			return
		}
	}

	return
}
