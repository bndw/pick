package main

func copyCommand(args ...string) int {
	safe, err := loadSafe()
	if err != nil {
		return handleError(err)
	}

	if err := safe.Copy(args[0]); err != nil {
		return handleError(err)
	}

	return 0
}
