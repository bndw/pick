package safe

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
)

func Load(password []byte, backend backends.Client, encryptionClient crypto.Client, config *config.Config) (*Safe, error) {
	s := Safe{
		backend:  backend,
		crypto:   encryptionClient,
		password: password,
		Config:   config,
	}

	data, err := s.backend.Load()
	if err != nil {
		if err == errors.ErrSafeNotFound {
			s.Accounts = make(map[string]Account)
			s.Notes = newNotesManager(&s)
			return &s, nil
		}
		return nil, err
	}

	safeDTO := safeDTO{}
	if err := json.Unmarshal(data, &safeDTO); err == nil { // nolint: vetshadow
		// Unmarshalling did succeed -> Safe uses new format
		// Do nothing
	} else {
		safeDTO.Ciphertext = data
		// If we don't have a config, use OpenPGP for backwards-compatibility
		defaultOpenPGPConfig := crypto.NewDefaultConfigWithType(crypto.ConfigTypeOpenPGP)
		safeDTO.Config = &defaultOpenPGPConfig
	}

	// We first try to decrypt a safe using its own configuration
	s.crypto, err = crypto.New(safeDTO.Config)
	if err != nil {
		return nil, err
	}

	upgradeCrypto := encryptionClient
	plaintext, err := s.crypto.Decrypt(safeDTO.Ciphertext, password)
	if err != nil {
		// Wasn't able to decrypt the safe with its own configuration.
		// Normally this shouldn't happen, never. As a fallback, we try
		// to decrypt using the user-provided config. If this fails, the
		// password is _definitely_ incorrect. If however decryption works
		// now, we need to upgrade the safe to use the user-provided config.
		// Therefore, store current crypto for possible upgrade.
		upgradeCrypto = s.crypto
		s.crypto = encryptionClient
		plaintext, err = s.crypto.Decrypt(safeDTO.Ciphertext, password)
		if err != nil {
			return nil, err
		}
	}

	var tmp Safe

	if err := json.Unmarshal(plaintext, &tmp); err != nil { // nolint: vetshadow
		return nil, errors.ErrSafeCorrupt
	}

	if tmp.Config != nil && tmp.Config.Version != "" {
		incompatible, err := versionIncompatible(tmp.Config.Version, config.Version) // nolint: vetshadow
		if err != nil {
			return nil, err
		} else if incompatible {
			return nil, fmt.Errorf("Safe is using a non-backwards compatible version. Please upgrade pick. pick version: %s, safe version: %s", config.Version, tmp.Config.Version)
		}
	}

	s.Accounts = tmp.Accounts
	s.Notes = newNotesManager(&s)
	if tmp.Notes != nil {
		s.Notes.Notes = tmp.Notes.Notes
	}

	// Default 'ModifiedOn' to 'CreatedOn' if it is empty
	for i, acc := range s.Accounts {
		if acc.ModifiedOn == 0 {
			acc.ModifiedOn = acc.CreatedOn
			s.Accounts[i] = acc
		}
	}

	// We need to compare the default / user-provided config with the safe config.
	// If they differ -> Upgrade safe
	if !reflect.DeepEqual(*safeDTO.Config, s.Config.Encryption) {
		fmt.Println("Upgrading safe")
		s.crypto = upgradeCrypto
		if err := s.save(); err != nil { // nolint: vetshadow
			fmt.Println("Error", err.Error())
		}
	}

	return &s, err
}

func versionIncompatible(v1, v2 string) (bool, error) {
	v1p, err := utils.ParseVersion(v1)
	if err != nil {
		return false, err
	}
	v2p, err := utils.ParseVersion(v2)
	if err != nil {
		return false, err
	}
	return v1p[0] > v2p[0] || v1p[1] > v2p[1], nil
}
