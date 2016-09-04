package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestCopy(t *testing.T) {
	testSafe := createTestSafe()
	defer removeTestSafe(testSafe)

	backend := backends.NewDiskBackend(&testSafe)
	safe, err := Load(testSafePassword, backend)
	if err != nil {
		t.Error(err)
	}

	if err = safe.Copy("foo"); err != nil {
		t.Error(err)
	}
}
