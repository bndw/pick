package backends

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/bndw/pick/errors"
)

const (
	defaultS3Bucket = "pick"
	defaultS3Key    = "safes/default.safe"
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
	Bucket  string
	Key     string
	Session *session.Session
	svc     *s3.S3
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
	if overrideAWSRegion && overrideAWSProfile {
		// Override both region and profile
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(region)},
			Profile: profile,
		}))
	} else if overrideAWSRegion {
		// Override only region
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))
	} else if overrideAWSProfile {
		// Override only profile
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
		}))
	} else {
		// Fallback to defaults/environment
		sess = session.Must(session.NewSession())
	}

	return &S3Backend{
		Bucket:  bucket,
		Key:     key,
		Session: sess,
		svc:     s3.New(sess),
	}, nil
}

func (s *S3Backend) Load() ([]byte, error) {
	result, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				// The safe at Key does not exist
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
				return nil, errors.ErrSafeNotFound
			default:
				// Some other AWS error
				fmt.Println(aerr.Error())
				return nil, err
			}
		}

		// A generic error
		fmt.Println(err.Error())
		return nil, err
	}

	defer result.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, result.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *S3Backend) Save(data []byte) error {
	_, err := s.svc.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			// TODO(bndw): Handle specific PutObject errors
			default:
				return err
			}
		}
		return err
	}

	return nil
}

func (s *S3Backend) Backup() error {
	return fmt.Errorf("S3 Backups are not implemented!")
}

func (s *S3Backend) SafeLocation() string {
	return fmt.Sprintf("s3://%s/%s", s.Bucket, s.Key)
}
