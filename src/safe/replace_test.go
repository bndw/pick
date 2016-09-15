package safe

import (
	"testing"
)

func TestReplace(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	account, err := safe.Replace("foo", "Bubbles", "kitt3ns")
	if account.Username != "Bubbles" {
		t.Errorf("Expected username Bubbles, got %s", account.Username)
	}
}
