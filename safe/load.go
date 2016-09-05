package safe

import (
	"encoding/json"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/errors"
)

func Load(password []byte, backend backends.Client, encryptionClient crypto.Client) (*Safe, error) {
	s := Safe{
		backend:  backend,
		crypto:   encryptionClient,
		password: password,
	}

	ciphertext, err := s.backend.Load()

	if _, ok := err.(*errors.SafeNotFound); ok {
		s.Accounts = make(map[string]Account)
		return &s, nil
	}
	if err != nil {
		return nil, err
	}

	plaintext, err := s.crypto.Decrypt(ciphertext, password)
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
