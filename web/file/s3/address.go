package s3

import (
	"context"
	"errors"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Error Messages during connect
var (
	ErrNoPassword = errors.New("no secret access key found")
)

// Connect try to use a path to setup a connection to s3 server
func Connect(path string) (*minio.Client, string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, "", err
	}

	tls := u.Scheme[len(u.Scheme)-1] == 's'
	accessKeyID := u.User.Username()
	secretAccessKey, ok := u.User.Password()
	if !ok {
		return nil, "", ErrNoPassword
	}
	query := u.Query()
	bucketName := query.Get("bucket")
	location := query.Get("location")

	u.User = nil
	u.RawQuery = ""

	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(u.String(), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: tls,
	})
	if err != nil {
		return nil, "", err
	}

	// create and check for bucket
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		if exists, err := minioClient.BucketExists(ctx, bucketName); err != nil || !exists {
			return nil, "", err
		}
	}

	return minioClient, bucketName, err
}
