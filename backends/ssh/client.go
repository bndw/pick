package ssh

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/errors"
)

type client struct {
	path         string
	backupConfig backends.BackupConfig
}

// Load loads data from a remote host over ssh
func (c *client) Load() ([]byte, error) {
	result, err := c.getObject(c.Bucket, c.Key)
	if err != nil {
		return nil, err
	}
	defer result.Close() // nolint: errcheck

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, result); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Save writes data to S3
func (c *client) Save(data []byte) error {
	return c.putObject(bytes.NewReader(data), c.Bucket, c.Key)
}

// SafeLocation return the path to the current safe
func (c *client) SafeLocation() string {
	return fmt.Sprintf("s3://%s/%s", c.Bucket, c.Key)
}

// Backup creates a backup of the current safe
func (c *client) Backup() error {
	var (
		now       = time.Now().Format(defaultS3BackupTimeFormat)
		backupKey = fmt.Sprintf(defaultS3BackupKeyTemplate, now)
	)

	data, err := c.Load()
	if err != nil {
		return err
	}

	return c.putObject(bytes.NewReader(data), c.Bucket, backupKey)
}

func (c *client) getObject(bucket, key string) (io.ReadCloser, error) {
	result, err := c.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				// The safe at Key does not exist
				return nil, errors.ErrSafeNotFound
			default:
				return nil, err
			}
		}

		return nil, err
	}

	return result.Body, nil
}

func (c *client) putObject(data io.Reader, bucket, key string) error {
	_, err := c.svc.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(data),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO(bndw): Handle specific PutObject errors
			default:
				return err
			}
		}
	}

	return err
}
