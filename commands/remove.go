package commands

import "fmt"

func RemoveCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	if err := safe.Remove(args[0]); err != nil {
		return handleError(err)
	}

	fmt.Println("Credential removed")
	return 0
}
