package s3

import (
	"github.com/bndw/pick/backends"
)

const (
	ClientName     = "s3"
	clientPriority = 0
)

func Register() {
	backends.Register(ClientName, clientPriority, _new)
}
