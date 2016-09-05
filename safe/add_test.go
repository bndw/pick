package safe

import (
	"testing"
)

func TestAdd(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}

	if err = safe.Add("github", "bndw", "fooBarBaz"); err != nil {
		t.Error(err)
	}
}
