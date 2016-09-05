package safe

import (
	"testing"
)

func TestCopy(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	if err = safe.Copy("foo"); err != nil {
		t.Error(err)
	}
}
