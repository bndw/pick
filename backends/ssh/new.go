package ssh

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bndw/pick/backends"
)

const (
	defaultSafePath = "/var/lib/pick/pick.safe"
	safeFileMode    = 0600
	safeDirMode     = 0700
)

// _new returns a new ssh client implementation.
func _new(config *backends.Config) (backends.Client, error) {
	safePath, ok := config.Settings["path"].(string)
	if !ok {
		safePath = defaultSafePath
	}

	config.Backup.DirPath = fmt.Sprintf(defaultBackupDir, homeDir, defaultSafeDirName)

	return &client{
		Bucket: bucket,
		Key:    key,
		svc:    s3.New(sess),
	}, nil
}
