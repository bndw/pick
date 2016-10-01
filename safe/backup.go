package safe

import (
	"github.com/bndw/pick/backends"
)

func Backup(client backends.Client) error {
	if err := client.Backup(); err != nil {
		return err
	}
	return nil
}
