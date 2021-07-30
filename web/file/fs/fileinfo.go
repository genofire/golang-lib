package fs

import (
	"os"

	"github.com/google/uuid"
)

type FileInfo struct {
	id          uuid.UUID
	name        string
	contentType string
	os.FileInfo
}

func (fi FileInfo) ID() uuid.UUID       { return fi.id }
func (fi FileInfo) ContentType() string { return fi.contentType }
func (fi FileInfo) Name() string        { return fi.name }
func (fi FileInfo) Sys() interface{}    { return fi.FileInfo }
