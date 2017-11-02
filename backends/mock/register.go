package mock

import (
	"github.com/bndw/pick/backends"
)

const (
	// Make sure to not export the mock client name
	clientName     = "_mock"
	clientPriority = 0
)

func init() {
	backends.Register(clientName, clientPriority, _new)
}
