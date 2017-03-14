package errors

import (
	"errors"
)

var (
	ErrSafeNotFound         = errors.New("Safe not found")
	ErrSafeCorrupt          = errors.New("Safe currupt")
	ErrSafeDecryptionFailed = errors.New("Unable to unlock safe with provided password")

	ErrBackupDisabled   = errors.New("Backups are disabled, increase `max_backups` to re-enable")
	ErrBackupFileExists = errors.New("Backup file already exists")

	ErrAccountAlreadyExists = errors.New("Account already exists")
	ErrAccountNotFound      = errors.New("Account not found")

	ErrInvalidCommandUsage = errors.New("Invalid use of command")
)
