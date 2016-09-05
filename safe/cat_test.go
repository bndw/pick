package safe

import (
	"testing"
)

func TestCat(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	if _, err = safe.Cat("foo"); err != nil {
		t.Error(err)
	}
}
