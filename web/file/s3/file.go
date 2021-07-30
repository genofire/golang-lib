package s3

import (
	"io/fs"

	"github.com/minio/minio-go/v7"
)

type File struct {
	*minio.Object
}

func (f File) Stat() (fs.FileInfo, error) {
	var fi FileInfo
	var err error
	fi.ObjectInfo, err = f.Object.Stat()
	if err != nil {
		return nil, err
	}
	return fi, nil
}
