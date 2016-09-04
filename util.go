package main

import (
	"fmt"

	"github.com/bndw/pick/safe"
	"github.com/bndw/pick/utils"
)

func loadSafe() (*safe.Safe, error) {
	password, err := utils.GetPasswordInput("Enter your master password")
	if err != nil {
		return nil, err
	}

	return safe.Load(password, nil)
}

func handleError(err error) int {
	fmt.Println(err)
	return 1
}
