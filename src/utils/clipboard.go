package utils

import (
	"github.com/atotto/clipboard"
)

func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}
