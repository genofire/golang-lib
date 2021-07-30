package file

import (
	"io/fs"

	"github.com/google/uuid"
)

type FileInfo interface {
	fs.FileInfo
	ID() uuid.UUID
	ContentType() string
}
