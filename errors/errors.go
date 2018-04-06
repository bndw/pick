package errors

import (
	"errors"
)

var (
	ErrSafeNotFound         = errors.New("Safe not found")
	ErrSafeCorrupt          = errors.New("Safe currupt")
	ErrSafeDecryptionFailed = errors.New("Unable to unlock safe with provided password")
	ErrSafeNotWritable      = errors.New("Safe not writable, this should never happen")

	ErrBackupDisabled   = errors.New("Backups are disabled, increase `max_backups` to re-enable")
	ErrBackupFileExists = errors.New("Backup file already exists")

	ErrAccountAlreadyExists = errors.New("Account already exists")
	ErrAccountNotFound      = errors.New("Account not found")

	ErrInvalidCommandUsage = errors.New("Invalid use of command")

	ErrAlreadyRunning = errors.New("Another instance of pick is already running. To prevent data loss, only a single pick-instance is allowed.")
)
