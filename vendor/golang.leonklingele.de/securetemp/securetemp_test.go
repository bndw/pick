package securetemp

import (
	"os"
	"testing"
)

const (
	size = DefaultSize
)

func TestTempDir(t *testing.T) {
	tmpDir, cleanupFunc, err := TempDir(size)
	if err != nil {
		t.Fatal("failed to create temp dir", err)
	}
	if _, err := os.Stat(tmpDir); err != nil {
		t.Fatal("something is wrong with the temp dir", err)
	}

	cleanupFunc()
	if _, err := os.Stat(tmpDir); err == nil {
		t.Fatal("temp dir should no longer exist after cleanup, but does", tmpDir)
	}
}

func TestTempFile(t *testing.T) {
	tmpFile, cleanupFunc, err := TempDir(size)
	if err != nil {
		t.Fatal("failed to create temp file", err)
	}
	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatal("something is wrong with the temp file", err)
	}

	cleanupFunc()
	if _, err := os.Stat(tmpFile); err == nil {
		t.Fatal("temp dir should no longer exist after cleanup, but does", tmpFile)
	}
}
