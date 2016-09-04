package commands

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
)

const (
	passwordLength = 25
)

func AddCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	name, username, password, errCode := parseAddArgs(args)
	if errCode > 0 {
		return errCode
	}

	err = safe.Add(name, username, password)
	if _, conflict := err.(*errors.AccountExists); conflict && overwrite(name) {
		rerr := safe.Remove(name)
		if rerr != nil {
			return handleError(rerr)
		}

		aerr := safe.Add(name, username, password)
		if aerr != nil {
			return handleError(aerr)
		}
	} else if err != nil {
		return handleError(err)
	}

	fmt.Println("Credential added")
	if utils.GetAnswer("Copy password to clipboard (Y/n)?", true) {
		if err := safe.Copy(name); err != nil {
			handleError(err)
		}
	}
	return 0
}

func overwrite(name string) bool {
	prompt := fmt.Sprintf("%s already exists, overwrite (y/N)?", name)
	return utils.GetAnswer(prompt, false)
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
		if utils.GetAnswer("Generate password (Y/n)?", true) {
			password, err = utils.GeneratePassword(passwordLength)
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
