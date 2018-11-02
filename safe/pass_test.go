package safe

import (
	"testing"

	"github.com/bndw/pick/config"
)

const (
	newPswd = "newPassword"
)

func TestChangePassword(t *testing.T) {
	safe, err := createTestSafe(t, true)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := safe.Add(accountName, initialUser, initialPswd); err != nil {
		t.Fatal(err)
	}

	if err := safe.ChangePassword([]byte(newPswd)); err != nil {
		t.Fatal(err)
	}

	config := &config.Config{
		Encryption: safe.Config.Encryption,
		Storage:    safe.Config.Storage,
		Version:    safe.Config.Version,
	}

	if _, err := Load(
		[]byte(newPswd),
		safe.backend,
		safe.crypto,
		config,
	); err != nil {
		t.Fatal(err)
	}
}
