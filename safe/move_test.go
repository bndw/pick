package safe

import (
	"testing"
)

func TestMove(t *testing.T) {
	safe, err := createTestSafe()
	if err != nil {
		t.Error(err)
	}
	defer removeTestSafe()

	name := "foo"
	newName := "foo-renamed"
	if err := safe.Move(name, newName); err != nil {
		t.Error(err)
	}
	if err := safe.Move(newName, name); err != nil {
		t.Error(err)
	}
}
