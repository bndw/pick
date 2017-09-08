package backends

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/bndw/pick/errors"
)

const (
	defaultS3Bucket = "pick"
	defaultS3Key    = "safes/default.safe"
	// Backups
	defaultS3BackupTimeFormat  = "2006-01-02_15-04-05"
	defaultS3BackupKeyTemplate = "backups/%s.safe"
)

// S3Backend provides AWS S3 backed safe storage
//
// This implementation relies on AWS Credentials in `~/.aws/credentials`
// for authentication. If AWS_PROFILE is not set, "default" will be used.
//
// The following environment variables are required:
// 		AWS_ACCESS_KEY_ID
//		AWS_SECRET_ACCESS_KEY
//		AWS_REGION
//
type S3Backend struct {
	Bucket string
	Key    string
	svc    *s3.S3
}

// NewS3Backend returns a new S3 Backend implementation.
//
// The following config settings are supported:
// 		region:  AWS Region to use
// 		profile: AWS Profile in ~/.aws/credentials to use
//		bucket:	 AWS S3 Bucket name for storing the safe. Defaults to `defaultS3Bucket`
//		key:	 	 AWS S3 Key name for storing the safe. Defaults to `defaultS3Key`
func NewS3Backend(config Config) (*S3Backend, error) {
	// AWS S3 Bucket overrides
	bucket, ok := config.Settings["bucket"].(string)
	if !ok {
		bucket = defaultS3Bucket
	}

	key, ok := config.Settings["key"].(string)
	if !ok {
		key = defaultS3Key
	}

	// TODO(bndw): Consider creating the bucket if it does not exist

	// AWS Session overrides
	region, overrideAWSRegion := config.Settings["region"].(string)
	profile, overrideAWSProfile := config.Settings["profile"].(string)

	var sess *session.Session
	switch {
	case overrideAWSRegion && overrideAWSProfile:
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(region)},
			Profile: profile,
		}))
	case overrideAWSRegion:
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
	case overrideAWSProfile:
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
		}))
	default:
		// Fallback to defaults/environment
		sess = session.Must(session.NewSession())
	}

	return &S3Backend{
		Bucket: bucket,
		Key:    key,
		svc:    s3.New(sess),
	}, nil
}

// Backup creates a backup of the current safe
func (s *S3Backend) Backup() error {
	var (
		now       = time.Now().Format(defaultS3BackupTimeFormat)
		backupKey = fmt.Sprintf(defaultS3BackupKeyTemplate, now)
	)

	data, err := s.Load()
	if err != nil {
		return err
	}

	return s.putObject(bytes.NewReader(data), s.Bucket, backupKey)
}

// Load loads data from S3
func (s *S3Backend) Load() ([]byte, error) {
	result, err := s.getObject(s.Bucket, s.Key)
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

// SafeLocation return the path to the current safe
func (s *S3Backend) SafeLocation() string {
	return fmt.Sprintf("s3://%s/%s", s.Bucket, s.Key)
}

// Save writes data to S3
func (s *S3Backend) Save(data []byte) error {
	return s.putObject(bytes.NewReader(data), s.Bucket, s.Key)
}

func (s *S3Backend) getObject(bucket, key string) (io.ReadCloser, error) {
	result, err := s.svc.GetObject(&s3.GetObjectInput{
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

func (s *S3Backend) putObject(data io.Reader, bucket, key string) error {
	_, err := s.svc.PutObject(&s3.PutObjectInput{
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
