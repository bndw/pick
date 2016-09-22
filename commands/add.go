package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "add [name] [username] [password]",
		Short: "Add a credential",
		Long: `The add command is used to add a new credential.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Add(args...))
		},
	})
}

func Add(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	name, username, password, errCode := parseAddArgs(args)
	if errCode > 0 {
		return errCode
	}

	account, err := safe.Add(name, username, password)
	if _, conflict := err.(*errors.AccountExists); conflict && overwrite(name) {
		var replaceErr error
		if account, replaceErr = safe.Replace(name, username, password); replaceErr != nil {
			return handleError(replaceErr)
		}
	} else if err != nil {
		return handleError(err)
	}

	fmt.Println("Credential added")
	if utils.Confirm("Copy password to clipboard", true) {
		if err := utils.CopyToClipboard(account.Password); err != nil {
			return handleError(err)
		}
	}
	return 0
}

func overwrite(name string) bool {
	prompt := fmt.Sprintf("%s already exists, overwrite", name)
	return utils.Confirm(prompt, false)
}

func parseAddArgs(args []string) (name, username, password string, errCode int) {
	if len(args) > 3 {
		fmt.Println("Usage: add [name] [username] [password]")
		return "", "", "", 1
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

	errCode = 1
	var err error

	if name == "" {
		if name, err = utils.GetInput("Enter a credential name"); err != nil {
			fmt.Println(err)
			return
		}
	}

	if username == "" {
		if username, err = utils.GetInput(fmt.Sprintf("Enter a username for %s", name)); err != nil {
			fmt.Println(err)
			return
		}
	}

	if password == "" {
		if utils.Confirm("Generate password", true) {
			password, err = utils.GeneratePassword(config.General.PasswordLen)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			var _password []byte
			if _password, err = utils.GetPasswordInput(fmt.Sprintf("Enter a password for %s", name)); err != nil {
				fmt.Println(err)
				return
			}

			password = string(_password)
		}
	}

	errCode = 0
	return
}
