package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "add [name] [username] [password]",
		Short: "Add a credential",
		Long:  "The add command is used to add a new credential.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Add, cmd, args)
		},
	})
}

func Add(args []string, flags *pflag.FlagSet) error {
	safe, err := loadSafe()
	if err != nil {
		return err
	}

	name, username, password, err := parseAddArgs(args)
	if err != nil {
		return err
	}

	account, err := safe.Add(name, username, password)
	if _, conflict := err.(*errors.AccountExists); conflict && overwrite(name) {
		var editErr error
		if account, editErr = safe.Edit(name, username, password); editErr != nil {
			return editErr
		}
	} else if err != nil {
		return err
	}

	fmt.Println("Credential added")
	if utils.Confirm("Copy password to clipboard", true) {
		if err := utils.CopyToClipboard(account.Password); err != nil {
			return err
		}
	}
	return nil
}

func overwrite(name string) bool {
	prompt := fmt.Sprintf("%s already exists, overwrite", name)
	return utils.Confirm(prompt, false)
}

func parseAddArgs(args []string) (name, username, password string, err error) {
	if len(args) > 3 {
		err = &errors.InvalidCommandUsage{}
		return
	}

	switch len(args) {
	case 3:
		password = args[2]
		fallthrough
	case 2:
		username = args[1]
		fallthrough
	case 1:
		name = args[0]
	}

	if name == "" {
		if name, err = utils.GetInput("Enter a credential name"); err != nil {
			return
		}
	}

	if username == "" {
		if username, err = utils.GetInput(fmt.Sprintf("Enter a username for %s", name)); err != nil {
			return
		}
	}

	if password == "" {
		if utils.Confirm("Generate password", true) {
			if password, err = utils.GeneratePassword(config.General.PasswordLen); err != nil {
				return
			}
		} else {
			var _password []byte
			if _password, err = utils.GetPasswordInput(fmt.Sprintf("Enter a password for %s", name)); err != nil {
				return
			}

			password = string(_password)
		}
	}

	return
}
