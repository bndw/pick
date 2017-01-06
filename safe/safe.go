package safe

import (
	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
)

type Safe struct {
	backend  backends.Client
	crypto   crypto.Client
	Config   *config.Config
	password []byte
	Accounts map[string]Account `json:"accounts"`
	Notes    *notesManager
}

type safeDTO struct {
	Config     *crypto.Config `json:"config"`
	Ciphertext []byte         `json:"ciphertext"`
}
