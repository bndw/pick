package safe

import (
	"testing"
)

func TestReplace(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}

	if err = safe.Replace("foo", "Bubbles", "kitt3ns"); err != nil {
		t.Error(err)
	}

	account, err := safe.Cat("foo")
	if err != nil {
		t.Error(err)
	}

	if account.Username != "Bubbles" {
		t.Errorf("Expected username Bubbles, got %s", account.Username)
	}
}
