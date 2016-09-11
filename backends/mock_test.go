package backends_test

import (
	"reflect"
	"testing"

	"github.com/bndw/pick/backends"
)

var (
	mockTestData = []byte("This is only a test")
)

func TestNewMockBackend(t *testing.T) {
	if b := backends.NewMockBackend(); b == nil {
		t.Errorf("Expected new mockBackend, got nil")
	}
}

func TestMockBackup(t *testing.T) {
	b := backends.NewMockBackend()

	if err := b.Backup(); err != nil {
		t.Errorf("Backup failed with %s", err)
	}
}

func TestMockLoad(t *testing.T) {
	b := backends.NewMockBackend()
	b.Data = mockTestData

	data, err := b.Load()
	if err != nil {
		t.Errorf("Load failed with %s", err)
	}
	if !reflect.DeepEqual(data, mockTestData) {
		t.Errorf("Load data returned unexpected data. Expected: '%s', Actual: '%s'", mockTestData, data)
	}
}

func TestMockSave(t *testing.T) {
	b := backends.NewMockBackend()

	if err := b.Save(mockTestData); err != nil {
		t.Errorf("Load failed with %s", err)
	}
}
