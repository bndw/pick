package clipboard

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	otherClipboard "github.com/atotto/clipboard"
)

const (
	clearClipboardCmd = "clear-clipboard"
	hashLen           = 4 // in bytes
)

// Copy copies the provided text to the system clipboard
// If clearAfter is given and larger 0, the clipboard is automatically
// cleared after the duration has passed
func Copy(text string, clearAfter Duration) error {
	if clearAfter.Seconds() > 0 {
		// Clear the clipboard after 'clearAfter'
		_ = launchClearer(text, clearAfter)
	}
	return otherClipboard.WriteAll(text)
}

// ClearIfMatch clears the clipboard if its current content's truncated
// hash equals 'match'
func ClearIfMatch(match string) error {
	// We're about to clear the clipboard
	// Ensure that we still have the expected contents in it
	// This avoids clearing other data which have been copied in the meantime
	current, err := otherClipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read current clipboard content: %s", err)
	}
	currentHashed := hashCurrentAsHex(current)
	if subtle.ConstantTimeCompare([]byte(currentHashed), []byte(match)) != 1 {
		return errors.New("not clearing clipboard as content has changed")
	}
	return Copy("", Duration{-1})
}

func launchClearer(current string, duration Duration) error {
	currentHashed := hashCurrentAsHex(current)

	pickPath, err := exec.LookPath("pick")
	if err != nil {
		return fmt.Errorf("failed to find 'pick' binary: %s", err)
	}
	secs := strconv.FormatInt(int64(duration.Seconds()), 10)
	cmd := exec.Command(pickPath, clearClipboardCmd, secs, currentHashed)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'pick': %s", err)
	}
	// Don't cmd.Wait() to run as background process
	return nil
}

func hashCurrentAsHex(current string) string {
	h := sha256.New()
	h.Write([]byte(current))
	return hex.EncodeToString(h.Sum(nil)[:hashLen])
}
