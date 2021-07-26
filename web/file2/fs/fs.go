/*
Package fs implements a non-hierarchical file store using the underlying (disk)
file system.

A file store contains directories named after UUIDs, each containing the
following files:
    name    the file name (``basename'')
    type    the file's MIME type
    data    the file data

This is not meant for anything for which the word ``scale'' plays any role at
all, ever, anywhere.

TODO: should io/fs ever support writable file systems, use one of those instead
of the ``disk'' (os.Open & Co.) (see https://github.com/golang/go/issues/45757)
*/
package fs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/google/uuid"
)

// FS is an on-disk file store.
type FS struct {
	// Root is the path of the file store, relative to the process's working
	// directory or absolute.
	Root string
}

func (fs *FS) Store(id uuid.UUID, name, contentType string, data io.Reader) error {
	p := path.Join(fs.Root, id.String())

	_, err := os.Stat(p)
	if err == nil {
		err = os.RemoveAll(p)
		if err != nil {
			return fmt.Errorf("remove old %v: %w", id, err)
		}
	}

	err = os.Mkdir(p, 0o750)
	if err != nil {
		return fmt.Errorf("mkdir %s: %w", id.String(), err)
	}

	defer func() {
		if err != nil {
			os.RemoveAll(p)
		}
	}()

	err = os.WriteFile(path.Join(p, "name"), []byte(name), 0o640)
	if err != nil {
		return fmt.Errorf("write name: %w", err)
	}
	err = os.WriteFile(path.Join(p, "type"), []byte(contentType), 0o640)
	if err != nil {
		return fmt.Errorf("write type: %w", err)
	}
	f, err := os.Create(path.Join(p, "data"))
	if err != nil {
		return fmt.Errorf("create data file: %w", err)
	}
	_, err = io.Copy(f, data)
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}
	return nil
}

func (fs *FS) RemoveUUID(id uuid.UUID) error {
	return os.RemoveAll(path.Join(fs.Root, id.String()))
}

// OpenUUID opens the file with the given UUID.
func (fs *FS) OpenUUID(id uuid.UUID) (fs.File, error) {
	f, err := os.Open(path.Join(fs.Root, id.String(), "data"))
	if err != nil {
		return nil, fmt.Errorf("open data: %w", err)
	}
	return File{File: f, fs: fs, id: uuid.MustParse(id.String())}, nil
}

// Open searches for and opens the file with the given name.
func (fs *FS) Open(name string) (fs.File, error) {
	files, err := os.ReadDir(fs.Root)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		id := v.Name()
		entryName, err := os.ReadFile(path.Join(fs.Root, id, "name"))
		if err != nil || string(entryName) != name {
			continue
		}
		f, err := os.Open(path.Join(fs.Root, id, "data"))
		if err != nil {
			return nil, fmt.Errorf("open data: %w", err)
		}
		return File{File: f, fs: fs, id: uuid.MustParse(id)}, nil
	}
	return nil, errors.New("no such file")
}

// Check checks whether the file store is operable, ie. whether fs.Root exists
// and is a directory.
func (fs *FS) Check() error {
	fi, err := os.Stat(fs.Root)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("file store is not a directory")
	}
	return nil
}
