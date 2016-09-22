package safe

import (
	"testing"
)

func TestGet(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	if _, err = safe.Get("foo"); err != nil {
		t.Error(err)
	}
}
