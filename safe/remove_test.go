package safe

import (
	"testing"
)

func TestRemove(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}

	if err = safe.Remove("foo"); err != nil {
		t.Error(err)
	}
}
