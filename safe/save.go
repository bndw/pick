package safe

import (
	"encoding/json"
)

func (s *Safe) save() error {
	plaintext, err := json.Marshal(s)
	if err != nil {
		return err
	}

	ciphertext, err := s.crypto.Encrypt(plaintext, s.password)
	if err != nil {
		return err
	}

	return s.backend.Save(ciphertext)
}
