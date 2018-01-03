package ssh

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestNewDefaultS3Backend(t *testing.T) {
	backend := newS3Backend(t, nil)
	if backend == nil {
		t.Fatalf("Expected new S3Backend, got nil")
	}

	if backend.Bucket != defaultS3Bucket {
		t.Errorf("Expected new bucket %s to be used, got %s", defaultS3Bucket, backend.Bucket)
	}
	if backend.Key != defaultS3Key {
		t.Errorf("Expected new key %s to be used, got %s", defaultS3Key, backend.Key)
	}
}

func TestNewS3BackendWithBucket(t *testing.T) {
	const bucket = "pick"

	backend := newS3Backend(t, &backends.Config{
		Settings: map[string]interface{}{
			"bucket": bucket,
		},
	})
	if backend == nil {
		t.Fatalf("Expected new S3Backend, got nil")
	}

	if backend.Bucket != bucket {
		t.Errorf("Expected new bucket %s to be used, got %s", bucket, backend.Bucket)
	}
	if backend.Key != defaultS3Key {
		t.Errorf("Expected new key %s to be used, got %s", defaultS3Key, backend.Key)
	}
}

func TestNewS3BackendWithBucketAndKey(t *testing.T) {
	const (
		bucket = "pick"
		key    = "public/default.safe"
	)

	backend := newS3Backend(t, &backends.Config{
		Settings: map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	})
	if backend == nil {
		t.Fatalf("Expected new S3Backend, got nil")
	}

	if backend.Bucket != bucket {
		t.Errorf("Expected new bucket %s to be used, got %s", bucket, backend.Bucket)
	}
	if backend.Key != key {
		t.Errorf("Expected new key %s to be used, got %s", key, backend.Key)
	}
}

func TestNewS3BackendWithAllOverrides(t *testing.T) {
	const (
		bucket  = "pick"
		key     = "public/foo.safe"
		region  = "us-west-2"
		profile = "dev"
	)

	backend := newS3Backend(t, &backends.Config{
		Settings: map[string]interface{}{
			"bucket":  bucket,
			"key":     key,
			"region":  region,
			"profile": profile,
		},
	})
	if backend == nil {
		t.Fatalf("Expected new S3Backend, got nil")
	}

	if backend.Bucket != bucket {
		t.Errorf("Expected new bucket %s to be used, got %s", bucket, backend.Bucket)
	}
	if backend.Key != key {
		t.Errorf("Expected new key %s to be used, got %s", key, backend.Key)
	}
}

func newS3Backend(t *testing.T, config *backends.Config) *client {
	// t.Helper() // TOOD(leon): Go 1.9 only :(
	if config == nil {
		tmp := backends.NewDefaultConfig()
		config = &tmp
	}
	c, err := backends.NewWithType(ClientName, config)
	if err != nil {
		t.Fatalf("Failed to create S3 backend: %v", err)
	}
	return c.(*client)
}

func init() {
	Register()
}
