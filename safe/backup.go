package safe

import (
	"github.com/bndw/pick/backends"
)

func Backup(client backends.Client) error {
	return client.Backup()
}
