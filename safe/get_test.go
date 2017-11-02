package safe

import (
	"testing"
)

func TestGet(t *testing.T) {
	safe, err := createTestSafe(t)
	if err != nil {
		t.Error(err)
	}

	if _, err = safe.Get("foo"); err != nil {
		t.Error(err)
	}
}
