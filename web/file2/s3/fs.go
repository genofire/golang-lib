/*
Package s3 implements a non-hierarchical file store using Amazon s3. A file
store uses a single bucket. Each file is an object with its UUID as the name.
The file name is stored in the user-defined object metadata x-amz-meta-filename.
*/
package s3

import (
	"context"
	"errors"
	"io"
	"io/fs"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FS struct {
	client *minio.Client
	bucket string
}

// New ``connects'' to an s3 endpoint and creates a file store using the
// specified bucket. The bucket is created if it doesn't exist.
func New(endpoint string, secure bool, id, secret, bucket, location string) (*FS, error) {
	var fs FS
	var err error
	fs.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(id, secret, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()

	err = fs.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region: location,
	})
	if err != nil {
		if exists, err := fs.client.BucketExists(ctx, bucket); err != nil || !exists {
			return nil, err
		}
	}
	fs.bucket = bucket

	return &fs, nil
}

func (fs *FS) Store(id uuid.UUID, name, contentType string, data io.Reader) error {
	ctx := context.TODO()
	_, err := fs.client.PutObject(ctx, fs.bucket, id.String(), data, -1, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"filename": name,
		},
		ContentType: contentType,
	})
	return err
}

func (fs *FS) RemoveUUID(id uuid.UUID) error {
	ctx := context.TODO()
	return fs.client.RemoveObject(ctx, fs.bucket, id.String(), minio.RemoveObjectOptions{})
}

// TODO: implement
func (fs *FS) Open(name string) (fs.File, error) {
	return nil, errors.New("not implemented")
}

func (fs *FS) OpenUUID(id uuid.UUID) (fs.File, error) {
	ctx := context.TODO()

	object, err := fs.client.GetObject(ctx, fs.bucket, id.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return File{Object: object}, nil
}

// TODO: do some checking
func (fs *FS) Check() error {
	return nil
}
