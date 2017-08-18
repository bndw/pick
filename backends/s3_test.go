package backends

import (
	"testing"

	"github.com/bndw/pick/backends"
)

func TestNewDefaultS3Backend(t *testing.T) {
	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
	})
	if err != nil {
		t.Fatal(err)
	}

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

	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
		Settings: map[string]interface{}{
			"bucket": bucket,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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

	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
		Settings: map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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

	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
		Settings: map[string]interface{}{
			"bucket":  bucket,
			"key":     key,
			"region":  region,
			"profile": profile,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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
