package commands

import (
	"fmt"

	"github.com/bndw/pick/utils"
)

func CatCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	account, err := safe.Cat(args[0])
	if err != nil {
		return handleError(err)
	}

	fmt.Printf(`account:  %s
username: %s
password: %s
created:  %s
`, account.Name, account.Username, account.Password,
		utils.FormatUnixTime(account.CreatedOn))
	return 0
}
