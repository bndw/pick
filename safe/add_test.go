package safe

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestAdd(t *testing.T) {
	testSafe := createTestSafe()
	defer removeTestSafe(testSafeName)

	backend := backends.NewDiskBackend(&testSafe)
	safe, err := Load(testSafePassword, backend)
	if err != nil {
		t.Error(err)
	}

	if err = safe.Add("github", "bndw", "fooBarBaz"); err != nil {
		t.Error(err)
	}
}
