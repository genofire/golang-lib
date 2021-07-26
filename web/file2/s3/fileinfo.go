package s3

import (
	"io/fs"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type FileInfo struct {
	minio.ObjectInfo
}

func (fi FileInfo) Name() string { return fi.UserMetadata["filename"] }
func (fi FileInfo) Size() int64  { return fi.ObjectInfo.Size }

// TODO: try to map s3 permissions to these, somehow
func (fi FileInfo) Mode() fs.FileMode   { return 0o640 }
func (fi FileInfo) ModTime() time.Time  { return fi.LastModified }
func (fi FileInfo) IsDir() bool         { return false }
func (fi FileInfo) Sys() interface{}    { return fi.ObjectInfo }
func (fi FileInfo) ID() uuid.UUID       { return uuid.MustParse(fi.Key) }
func (fi FileInfo) ContentType() string { return fi.ObjectInfo.ContentType }
