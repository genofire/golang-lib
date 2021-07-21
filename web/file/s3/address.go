package s3

import (
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Connect try to use a path to setup a connection to s3 server
func Connect(path string) (client, string, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, "", err
	}

	secretAccessKey, err := url.Userinfo.Password()
	if err != nil {
		return nil, "", err
	}

	// Initialize minio client object.
	minioClient, err := minio.New(url.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(url.Userinfo.Username(), secretAccessKey, ""),
		Secure: url.Schema[-1] == "s",
	})
	return minioClient, url.Path, err
}
