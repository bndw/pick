package main

import (
	"fmt"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
)

const (
	passwordLength = 25
)

func addCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	name, username, password, errCode := parseAddArgs(args)
	if errCode > 0 {
		return errCode
	}

	err = safe.Add(name, username, password)
	if _, conflict := err.(*errors.AccountExists); conflict {
		if overwrite(name) {
			// Remove the existing credential
			rerr := safe.Remove(name)
			if rerr != nil {
				return handleError(rerr)
			}

			addCommand([]string{name, username, password}...)
		}
	} else if err != nil {
		return handleError(err)
	}

	fmt.Println("Credential added")
	return 0
}

func overwrite(name string) bool {
	prompt := fmt.Sprintf("%s already exists, overwrite (y/n) ?\n", name)
	return utils.GetAnswer(prompt)
}

func parseAddArgs(args []string) (name, username, password string, errCode int) {
	if len(args) > 3 {
		fmt.Println("Usage: add [name] [username] [password]")
		return "", "", "", 1
	}

	switch len(args) {
	case 3:
		name = args[0]
		username = args[1]
		password = args[2]
	case 2:
		name = args[0]
		username = args[1]
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
		if utils.GetAnswer("Generate password (y/n)?") {
			password = utils.GeneratePassword(passwordLength)
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
