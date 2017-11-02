package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bndw/pick/backends"
)

const (
	defaultS3Bucket = "pick"
	defaultS3Key    = "safes/default.safe"
)

// _new returns a new S3 client implementation.
//
// The following config settings are supported:
// 		region:  AWS Region to use
// 		profile: AWS Profile in ~/.aws/credentials to use
//		bucket:	 AWS S3 Bucket name for storing the safe. Defaults to `defaultS3Bucket`
//		key:	 	 AWS S3 Key name for storing the safe. Defaults to `defaultS3Key`
func _new(config *backends.Config) (backends.Client, error) {
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

	return &client{
		Bucket: bucket,
		Key:    key,
		svc:    s3.New(sess),
	}, nil
}
