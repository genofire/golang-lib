/*
Package file abstracts non-hierarchical file stores. Each file consists of a
name, a MIME type, a UUID, and data. File names may be duplicate.

TODO: think about name vs. UUID again—should the ``name'' really be the filename
and not the UUID?
*/
package file

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// An FS is a file store.
type FS interface {
	// Store stores a new file with the given UUID, name, and MIME type.
	// Its data is taken from the provided Reader. If it encounters an
	// error, it does nothing. Any existing file with the same UUID is
	// overwritten.
	Store(id uuid.UUID, name, contentType string, data io.Reader) error
	// RemoveUUID deletes a file.
	RemoveUUID(id uuid.UUID) error
	// Open opens a file by its name. If multiple files have the same name,
	// it is unspecified which one is opened. This may very well be very
	// slow. This is bad. Go away.
	Open(name string) (fs.File, error)
	// OpenUUID opens a file by its UUID.
	OpenUUID(id uuid.UUID) (fs.File, error)
	// Check checks the health of the file store. If the file store is not
	// healthy, it returns a descriptive error. Otherwise, the file store
	// should be usable.
	Check() error
}

// StoreFromHTTP stores the first file with given form key from an HTTP
// multipart/form-data request. Its Content-Type header is ignored; the type is
// detected. The file name is the last part of the provided file name not
// containing any slashes or backslashes.
//
// TODO: store all files with the given key instead of just the first one
func StoreFromHTTP(fs FS, r *http.Request, key string) error {
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return fmt.Errorf("get file from request: %w", err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read from file: %w", err)
	}
	contentType := http.DetectContentType(buf[:n])
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek in file: %w", err)
	}

	i := strings.LastIndexAny(fileHeader.Filename, "/\\")
	// if i == -1 { i = -1 }

	return fs.Store(uuid.New(), fileHeader.Filename[i+1:], contentType, file)
}
