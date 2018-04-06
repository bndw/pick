package safe

import (
	"testing"
)

func TestList(t *testing.T) {
	safe, err := createTestSafe(t, false)
	if err != nil {
		t.Error(err)
	}

	accounts := safe.List()
	if len(accounts) < 1 {
		t.Error(err)
	}
}
