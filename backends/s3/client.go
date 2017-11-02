package s3

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bndw/pick/errors"
)

const (
	// Backups
	defaultS3BackupTimeFormat  = "2006-01-02_15-04-05"
	defaultS3BackupKeyTemplate = "backups/%s.safe"
)

// client provides AWS S3 backed safe storage
//
// This implementation relies on AWS Credentials in `~/.aws/credentials`
// for authentication. If AWS_PROFILE is not set, "default" will be used.
//
// The following environment variables are required:
// 		AWS_ACCESS_KEY_ID
//		AWS_SECRET_ACCESS_KEY
//		AWS_REGION
//
type client struct {
	Bucket string
	Key    string
	svc    *s3.S3
}

// Load loads data from S3
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
