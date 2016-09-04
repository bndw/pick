package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestList(t *testing.T) {
	testSafe := createTestSafe()
	defer removeTestSafe(testSafe)

	backend := backends.NewDiskBackend(&testSafe)
	safe, err := Load(testSafePassword, backend)
	if err != nil {
		t.Error(err)
	}

	accounts, err := safe.List()
	if err != nil {
		t.Error(err)
	}

	if len(accounts) < 1 {
		t.Error(err)
	}
}
