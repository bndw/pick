package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
	"github.com/bndw/pick/utils/pswdgen"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   `edit name ("username" | "password")`,
		Short: "Edit a credential",
		Long:  "The edit command is used to edit an existing credential.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Edit, cmd, args)
		},
	})
}

func Edit(args []string, flags *pflag.FlagSet) error {
	safe, err := loadSafe()
	if err != nil {
		return err
	}

	name, username, password, err := parseEditArgs(args)
	if err != nil {
		return err
	}

	account, err := safe.Edit(name, username, password)
	if err != nil {
		return err
	}

	fmt.Println("Credential updated")
	if utils.Confirm("Copy password to clipboard", true) {
		if err := utils.CopyToClipboard(account.Password); err != nil {
			return err
		}
	}

	return nil
}

func parseEditArgs(args []string) (name, username, password string, err error) {
	if len(args) > 2 {
		err = &errors.InvalidCommandUsage{}
		return
	}

	var action string
	switch len(args) {
	case 2:
		action = args[1]
		fallthrough
	case 1:
		name = args[0]
	}

	if action == "username" {
		if username, err = utils.GetInput(fmt.Sprintf("Enter a new username for %s", name)); err != nil {
			return
		}
	} else if action == "password" {
		if utils.Confirm("Generate new password", true) {
			password, err = pswdgen.Generate(config.General.Password)
			if err != nil {
				return
			}
		} else {
			var _password []byte
			if _password, err = utils.GetPasswordInput(fmt.Sprintf("Enter a new password for %s", name)); err != nil {
				return
			}

			password = string(_password)
		}
	} else {
		err = fmt.Errorf("Invalid edit action specified: %s\n", action)
		return
	}

	return
}
