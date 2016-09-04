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

	accounts := safe.List()
	if len(accounts) < 1 {
		fmt.Println("No accounts to export")
		return 1
	}

	utils.PrettyPrint(accounts)
	return 0
}
