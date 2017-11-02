package safe

import (
	"testing"
)

func TestEdit(t *testing.T) {
	safe, err := createTestSafe(t)
	if err != nil {
		t.Error(err)
	}

	account, err := safe.Edit("foo", "Bubbles", "kitt3ns")
	if err != nil {
		t.Error(err)
	}
	if account.Username != "Bubbles" {
		t.Errorf("Expected username Bubbles, got %s", account.Username)
	}
}
