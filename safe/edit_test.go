package safe

import (
	"testing"
)

func TestEdit(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	account, err := safe.Edit("foo", "Bubbles", "kitt3ns")
	if account.Username != "Bubbles" {
		t.Errorf("Expected username Bubbles, got %s", account.Username)
	}
}
