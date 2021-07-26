package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/google/uuid"
)

// File is a file handle with its associated ID.
type File struct {
	*os.File
	fs *FS
	id uuid.UUID
}

func (f File) Stat() (fs.FileInfo, error) {
	fi := FileInfo{id: f.id}
	var err error
	fi.FileInfo, err = f.File.Stat()
	if err != nil {
		return nil, fmt.Errorf("os stat: %w", err)
	}
	name, err := os.ReadFile(path.Join(f.fs.Root, f.id.String(), "name"))
	if err != nil {
		return nil, fmt.Errorf("reading name: %w", err)
	}
	fi.name = string(name)
	contentType, err := os.ReadFile(path.Join(f.fs.Root, f.id.String(), "type"))
	if err != nil {
		return nil, fmt.Errorf("reading type: %w", err)
	}
	fi.contentType = string(contentType)
	return fi, nil
}
