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

type SafeDecryptionFailed struct {
}

func (e *SafeDecryptionFailed) Error() string {
	return "Unable to unlock safe with provided password"
}

type BackupDisabled struct {
}

func (e *BackupDisabled) Error() string {
	return "Backups are disabled, increase `max_backups` to re-enable"
}

type BackupFileExists struct {
}

func (e *BackupFileExists) Error() string {
	return "Backup file already exists"
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
