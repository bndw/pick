package main

import (
	"fmt"

	"github.com/bndw/pick/utils"
)

func exportCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	accounts, err := safe.List()
	if err != nil {
		return handleError(err)
	}

	if len(accounts) < 1 {
		fmt.Println("No accounts to export")
		return 1
	}

	var accountSlice []interface{}
	for _, account := range accounts {
		accountSlice = append(accountSlice, account)
	}

	utils.PrettyPrint(accountSlice)
	return 0
}
