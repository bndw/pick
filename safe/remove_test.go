package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestRemove(t *testing.T) {
	testSafe := createTestSafe()
	defer removeTestSafe(testSafe)

	backend := backends.NewDiskBackend(&testSafe)
	safe, err := Load(testSafePassword, backend)
	if err != nil {
		t.Error(err)
	}

	if err = safe.Remove("foo"); err != nil {
		t.Error(err)
	}
}
