package main

import (
	"fmt"
	"sort"
)

func listCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	var accountNames []string
	for _, account := range safe.List() {
		accountNames = append(accountNames, account.Name)
	}

	if len(accountNames) > 0 {
		sort.Strings(accountNames)
		for _, name := range accountNames {
			fmt.Println(name)
		}
	} else {
		fmt.Println("No accounts found")
	}

	return 0
}
