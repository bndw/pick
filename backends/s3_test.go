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

func TestS3BackendLoad(t *testing.T) {
	const (
		bucket = "bndw-pick"
		key    = "public/default.safe"
		region = "us-west-2"
	)

	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
		Settings: map[string]interface{}{
			"bucket": bucket,
			"key":    key,
			"region": region,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := backend.Load(); err != nil {
		t.Fatal(err)
	}
}

func TestS3BackendSave(t *testing.T) {
	const (
		bucket = "bndw-pick"
		key    = "public/default.safe"
		region = "us-west-2"
		data   = `{"config":{"type":"chachapoly","chachapoly":{"keyderivation":"pbkdf2","pbkdf2":{"hash":"sha512","iterations":100000,"saltlen":16}}},"ciphertext":"eyJzYWx0IjoidTdyYVMxRWFuL2hzNDlIVVMzV2g4QT09Iiwibm9uY2UiOiIzOFJhQmxldzgrSzlpOEdqIiwiY2lwaGVydGV4dCI6IlpUZnJrSlZOVjdLQW50Vi9wYTdhczFuUE1lQmYxMDRjRURCM29pRmtiN1JMTURyVGo5N29zdkFLWUxOa0lNRnliRVE2akFYMFJSK0NBV3pGLzE4UWtWb2hEeVdrV1BPNkdPVXFrTG5YamhTU0x4aTQrOWVmbHYzMll0dmNRZ2dEbXlnMkYxUXppdllONlFIaUQ4clBwMDZvWVYvdXNndGJEZG55THhrVnU0R1VxOFFLdTRyVlMrRWtOSndLd2hkcmx4YkV6VVNXZlphVzhHTkVBcysvWkxQR00xUC9MT3lLM0ZHQ1RJMmdmWUhMQU01eDNXTXZZcEM3bHlLNTBZS1JvaUlhcDJGZExuM1ZTSkxJZ2VtZkttNEp2ZWZlNWNBK0F2S2t6V05aUHJPL3ZqSWhKYUhTY0FsS3owVTdJVHZkZnI1aGZEc1BjcDFCc2RIRWN2SWs2TE9Va29sd3NBSkp1ZWpWeHBiRkhIVVBXSnRwanBlVGtHWEVNVXIxYjlKaHFpT0NDVGtGczRuam12QXJ6TDV1aGE4L1FMWnZEeitkZy9hZTAzbDUrRE9wSWFTaTRWYjVtQkRSM2RlbkhrUjNEUyttU3VXTGhxSUlwaDhLUlM4ejI3SFlqQVExc0duK2cvM0F6dGxZbVhCWk5tRHhXc3ZkU0w1cjZuaGtEUU9IY1d4MTUzUUlxZHd0ZWZFMDR0OVZJWnU2bGlsVDhRMklSV3cvUWZaR2UrNDNxQmFzcFFqYXFmaGRzT1lmV1pDVTM5VndCUTlZUzk2UjA1TllDcURPdE0rNlhzTnVjZGRZWnk5Y0Uxb0piSFplUXBLc2JoZHNESVVta04yR2JLZ3o2a01FSmlHNUZzQkQrQmJhcTd5akdocGVvYmtQd1BEeS9uQVQ5cnFXTjJEaHZIVGF4SCtyU2I1MWo0bjltNjJJalZlbVJSOGxWcDY0U21wSmkzMD0ifQ=="}`
	)

	backend, err := backends.NewS3Backend(backends.Config{
		Type: backends.ConfigTypeS3,
		Settings: map[string]interface{}{
			"bucket": bucket,
			"key":    key,
			"region": region,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Save([]byte(data)); err != nil {
		t.Fatal(err)
	}
}
