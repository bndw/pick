package safe

import (
	"encoding/json"
	"fmt"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/errors"
)

// Load loads, decrypts, and returns an instantiated Safe.
// In the case that the Safe configuration has changed, the Safe will be upgraded
// to the latest configuration and saved.
func Load(password []byte, backend backends.Client, encryptionClient crypto.Client, config *config.Config) (*Safe, error) {
	ciphertext, err := backend.Load()
	if err != nil {
		switch err {
		case errors.ErrSafeNotFound:
			return NewSafe(password, backend, encryptionClient, config, nil, nil)
		default:
			return nil, err
		}
	}

	var configurationHasChanged bool
	plaintext, err := decryptSafeDTO(ciphertext, password, nil)
	if err != nil {
		// Failed to decrypt the safe with its own configuration,
		// fall back to the user-provided config.
		plaintext, err = decryptSafeDTO(ciphertext, password, encryptionClient)
		if err != nil {
			// The password is definitely incorrect
			return nil, err
		}

		configurationHasChanged = true
	}

	var tmp Safe
	if err := json.Unmarshal(plaintext, &tmp); err != nil { // nolint: vetshadow
		return nil, errors.ErrSafeCorrupt
	}
	if tmp.Config != nil && tmp.Config.Version != "" {
		if err := tmp.IsVersionCompatible(config.Version); err != nil { // nolint: vetshadow
			return nil, err
		}
	}

	var notes map[string]note
	if tmp.Notes != nil {
		notes = tmp.Notes.Notes
	}

	s, err := NewSafe(password, backend, encryptionClient, config, tmp.Accounts, notes)
	if err != nil {
		return nil, err
	}

	if configurationHasChanged {
		if err := s.save(); err != nil { // nolint: vetshadow
			fmt.Println("Error saving file", err.Error())
		}
	}

	return s, err
}

// decryptSafeDTO attempts to decrypt the safeDTO ciphertext with the provided
// password and config
func decryptSafeDTO(safeBytes, password []byte, cryptoClient crypto.Client) ([]byte, error) {
	dto := safeDTO{}
	if err := json.Unmarshal(safeBytes, &dto); err != nil { // nolint: vetshadow
		// Safe likely uses a legacy format
		dto.Ciphertext = safeBytes
		// If we don't have a config, use OpenPGP for backwards-compatibility
		defaultOpenPGPConfig := crypto.NewDefaultConfigWithType(crypto.ConfigTypeOpenPGP)
		dto.Config = &defaultOpenPGPConfig
	}

	if cryptoClient == nil {
		var err error
		cryptoClient, err = crypto.New(dto.Config)
		if err != nil {
			return nil, err
		}
	}

	return cryptoClient.Decrypt(dto.Ciphertext, password)
}
