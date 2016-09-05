package safe

import (
	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
)

type Safe struct {
	backend  backends.Client
	crypto   crypto.Client
	password []byte
	Accounts map[string]Account `json:"accounts"`
}
