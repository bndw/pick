package commands

import (
	"fmt"
	"os"

	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   `edit name ("username" | "password")`,
		Short: "Edit a credential",
		Long: `The edit command is used to edit an existing credential.
            `,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(Edit(args...))
		},
	})
}

func Edit(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	name, username, password, errCode := parseEditArgs(args)
	if errCode > 0 {
		return errCode
	}

	account, err := safe.Edit(name, username, password)
	if err != nil {
		return handleError(err)
	}

	fmt.Println("Credential updated")
	if utils.Confirm("Copy password to clipboard", true) {
		if err := utils.CopyToClipboard(account.Password); err != nil {
			return handleError(err)
		}
	}

	return 0
}

func parseEditArgs(args []string) (name, username, password string, errCode int) {
	if len(args) > 2 {
		fmt.Println(`Usage: edit name ("username" | "password")`)
		return "", "", "", 1
	}

	var action string
	switch len(args) {
	case 2:
		action = args[1]
		fallthrough
	case 1:
		name = args[0]
	}

	errCode = 1
	var err error

	if action == "username" {
		if username, err = utils.GetInput(fmt.Sprintf("Enter a new username for %s", name)); err != nil {
			fmt.Println(err)
			return
		}
	} else if action == "password" {
		if utils.Confirm("Generate new password", true) {
			password, err = utils.GeneratePassword(config.General.PasswordLen)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			var _password []byte
			if _password, err = utils.GetPasswordInput(fmt.Sprintf("Enter a new password for %s", name)); err != nil {
				fmt.Println(err)
				return
			}

			password = string(_password)
		}
	} else {
		fmt.Printf("Invalid edit action specified: %s\n", action)
		return
	}

	errCode = 0
	return
}
