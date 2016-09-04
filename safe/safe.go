package safe

import (
	"github.com/bndw/pick/backends"
)

type Safe struct {
	backend  backends.Backend   `json:"-"`
	password []byte             `json:"-"`
	Accounts map[string]Account `json:"accounts"`
}
