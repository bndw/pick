package errors

type SafeNotFound struct {
}

func (e *SafeNotFound) Error() string {
	return "Safe not found"
}

type SafeCorrupt struct {
}

func (e *SafeCorrupt) Error() string {
	return "Safe currupt"
}

type AccountExists struct {
}

func (e *AccountExists) Error() string {
	return "Account exists"
}

type AccountNotFound struct {
}

func (e *AccountNotFound) Error() string {
	return "Account not found"
}
