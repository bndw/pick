package safe

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/errors"
)

// Load loads, decrypts, and returns an instantiated Safe.
// In the case that the Safe configuration has changed, the Safe will be upgraded
// to the latest configuration and saved.
func Load(password []byte, backendClient backends.Client, cryptoClient crypto.Client, config *config.Config) (*Safe, error) {
	ciphertext, err := backendClient.Load()
	if err != nil {
		switch err {
		case errors.ErrSafeNotFound:
			return New(password, backendClient, cryptoClient, config, nil, nil)
		default:
			return nil, err
		}
	}

	dto := NewSafeDTO(ciphertext)

	var configChanged bool
	plaintext, err := dto.Decrypt(password)
	if err != nil {
		// Failed to decrypt the safe with its own configuration, fall back to the
		// user-provided config.
		plaintext, err = dto.DecryptWithClient(password, cryptoClient)
		if err != nil {
			// The password is definitely incorrect
			return nil, err
		}
		configChanged = true
	} else if !reflect.DeepEqual(*dto.Config, config.Encryption) {
		// Although the safe was decrypted with its own configuration, the user's
		// safe config has changed.
		configChanged = true
	}

	var tmp Safe
	if err := json.Unmarshal(plaintext, &tmp); err != nil { // nolint: vetshadow
		return nil, errors.ErrSafeCorrupt
	}
	if tmp.Config != nil && tmp.Config.Version != "" {
		if err := tmp.RequireCompatibilityWith(config.Version); err != nil { // nolint: vetshadow
			return nil, err
		}
	}

	var notes map[string]note
	if tmp.Notes != nil {
		notes = tmp.Notes.Notes
	}

	s, err := New(password, backendClient, cryptoClient, config, tmp.Accounts, notes)
	if err != nil {
		return nil, err
	}

	if configChanged {
		w := backendClient.IsWritable()
		// Make safe writable
		backendClient.SetWritable(true)
		// Restore old mode
		defer backendClient.SetWritable(w)

		fmt.Println("Upgrading safe")
		if err := s.save(); err != nil { // nolint: vetshadow
			fmt.Println("Error saving safe", err.Error())
		}
	}

	return s, err
}
