package safe

import (
	"encoding/json"

	"github.com/bndw/pick/utils"
)

func (s *Safe) save() error {
	plaintext, err := json.Marshal(s)
	if err != nil {
		return err
	}

	ciphertext, err := utils.Encrypt(plaintext, _password)
	if err != nil {
		return err
	}

	return s.backend.Save(ciphertext)
}
