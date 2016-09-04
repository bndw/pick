package safe

import (
	"encoding/json"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
)

func Load(password []byte, backend backends.Backend) (*Safe, error) {
	if backend == nil {
		// TODO: Check the config for a disk location
		backend = backends.NewDiskBackend(nil)
	}
	s := Safe{backend: backend, password: password}

	ciphertext, err := s.backend.Load()

	if _, ok := err.(*errors.SafeNotFound); ok {
		s.Accounts = make(map[string]Account)
		return &s, nil
	}
	if err != nil {
		return nil, err
	}

	plaintext, err := utils.Decrypt(ciphertext, password)
	if err != nil {
		return nil, err
	}

	var tmp Safe

	if err := json.Unmarshal(plaintext, &tmp); err != nil {
		return nil, &errors.SafeCorrupt{}
	}

	s.Accounts = tmp.Accounts

	return &s, err
}
