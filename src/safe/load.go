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

func Load(password []byte, backend backends.Client, encryptionClient crypto.Client, config *config.Config) (*Safe, error) {
	s := Safe{
		backend:  backend,
		crypto:   encryptionClient,
		password: password,
		Config:   config,
	}

	data, err := s.backend.Load()
	if err != nil {
		if _, ok := err.(*errors.SafeNotFound); ok {
			s.Accounts = make(map[string]Account)
			return &s, nil
		}
		return nil, err
	}

	safeDTO := safeDTO{}
	if err := json.Unmarshal(data, &safeDTO); err == nil {
		// Unmarshalling did succeed -> Safe uses new format
		// Do nothing
	} else {
		safeDTO.Ciphertext = data
		// If we don't have a config, use OpenPGP for backwards-compatibility
		defaultOpenPGPConfig := crypto.NewDefaultConfigWithType(crypto.ConfigTypeOpenPGP)
		safeDTO.Config = &defaultOpenPGPConfig
	}

	upgradeNeeded := false
	plaintext, err := s.crypto.Decrypt(safeDTO.Ciphertext, password)
	if err != nil {
		// Wasn't able to decrypt the safe with the default / user-provided config.
		// Now use config from the safe. If this doesn't work, then the password
		// is _definitely_ incorrect. If however decryption works now, we need
		// to upgrade the safe to use the default / user-provided config.
		userCrypto := s.crypto
		s.crypto, err = crypto.New(safeDTO.Config)
		if err != nil {
			return nil, err
		}
		plaintext, err = s.crypto.Decrypt(safeDTO.Ciphertext, password)
		if err != nil {
			return nil, err
		}
		// Safe upgrade is needed, restore default / user-provided config
		upgradeNeeded = true
		s.crypto = userCrypto
	}

	var tmp Safe

	if err := json.Unmarshal(plaintext, &tmp); err != nil {
		return nil, &errors.SafeCorrupt{}
	}

	s.Accounts = tmp.Accounts

	// We still need to compare the default / user-provided config with the safe config.
	// If they differ -> Upgrade safe
	// This check is required for the OpenPGP mode, as it gets its config from the ciphertext.
	if upgradeNeeded || !reflect.DeepEqual(*safeDTO.Config, s.Config.Encryption) {
		fmt.Println("Upgrading safe")
		if err := s.save(); err != nil {
			fmt.Println("Error", err.Error())
		}
	}

	return &s, err
}
