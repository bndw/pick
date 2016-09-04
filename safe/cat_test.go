package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestCat(t *testing.T) {
	testSafe := createTestSafe()
	defer removeTestSafe(testSafe)

	backend := backends.NewDiskBackend(&testSafe)
	safe, err := Load(testSafePassword, backend)
	if err != nil {
		t.Error(err)
	}

	if _, err = safe.Cat("foo"); err != nil {
		t.Error(err)
	}
}
