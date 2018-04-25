package file

import (
	"github.com/bndw/pick/backends"
)

const (
	ClientName     = "file"
	clientPriority = 1
)

func Register() {
	backends.Register(ClientName, clientPriority, _new)
}
