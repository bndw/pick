package safe

import (
	"fmt"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/config"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/utils"
)

// Safe holds the structure for defining a Safe object
type Safe struct {
	backend  backends.Client
	crypto   crypto.Client
	Config   *config.Config
	password []byte
	Accounts map[string]Account `json:"accounts"`
	Notes    *notesManager
}

// New returns a new, initialized Safe
func New(password []byte, backend backends.Client, encryptionClient crypto.Client, config *config.Config, accounts map[string]Account, notes map[string]note) (*Safe, error) {
	s := Safe{
		backend:  backend,
		crypto:   encryptionClient,
		password: password,
		Accounts: accounts,
		Config:   config,
	}

	// Setup accounts
	if s.Accounts == nil {
		s.Accounts = make(map[string]Account)
	} else {
		// Default 'ModifiedOn' to 'CreatedOn' if it is empty
		for i, acc := range s.Accounts {
			if acc.ModifiedOn == 0 {
				acc.ModifiedOn = acc.CreatedOn
				s.Accounts[i] = acc
			}
		}
	}

	// Setup notes
	s.Notes = newNotesManager(&s)
	if notes != nil {
		s.Notes.Notes = notes
	}

	return &s, nil
}

// RequireCompatibilityWith returns an error if the Safe version is incompatible
// with the provided version
func (s *Safe) RequireCompatibilityWith(version string) error {
	versionIncompatible := func(v1, v2 string) (bool, error) {
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

	if s.Config == nil && s.Config.Version == "" {
		// No versioning, this is not an error
		return nil
	}

	incompatible, err := versionIncompatible(s.Config.Version, version)
	if err != nil {
		return err
	} else if incompatible {
		return fmt.Errorf("Safe is using a non-backwards compatible version. Please upgrade pick. pick version: %s, safe version: %s", version, s.Config.Version)
	} else {
		return nil
	}
}

// safeDTO holds the structure of how the Safe is persisted by the Backend
type safeDTO struct {
	Config     *crypto.Config `json:"config"`
	Ciphertext []byte         `json:"ciphertext"`
}
