package safe

import (
	"testing"
)

const (
	name     = "zombocom"
	username = "sarah"
	password = "welcome"
)

func TestAdd(t *testing.T) {
	safe, err := createTestSafe(t, true)
	if err != nil {
		t.Error(err)
	}

	account, err := safe.Add(name, username, password)
	if err != nil {
		t.Error(err)
	}

	if account.Username != username {
		t.Error("Unexpected account username:", username)
	}
	if account.Password != password {
		t.Error("Unexpected account password:", password)
	}
}
