package fs

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/google/uuid"

	"dev.sum7.eu/genofire/golang-lib/web/file"
)

// consts for filemanager
const (
	StorageTypeFS = "fs"
)

// error messages
var (
	ErrPathNotExistsOrNoDirectory = errors.New("path invalid: not exists or not an directory")
)

// FileManager to handle data on disk
type FileManager struct {
}

// Check if filemanager could be used
func (m *FileManager) Check(s *file.Service) error {
	info, err := os.Stat(s.Path)
	if os.IsNotExist(err) || !info.IsDir() {
		return ErrPathNotExistsOrNoDirectory
	}
	return nil
}

// Save a file on disk and update file db
func (m *FileManager) Save(s *file.Service, file *file.File, src multipart.File) error {
	file.ID = uuid.New()
	file.Path = path.Join(file.ID.String(), file.Filename)

	directory := path.Join(s.Path, file.ID.String())
	os.Mkdir(directory, 0750)

	out, err := os.Create(path.Join(s.Path, file.Path))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// Read get an reader of an file
func (m *FileManager) Read(s *file.Service, file *file.File) (io.Reader, error) {
	return os.Open(path.Join(s.Path, file.Path))
}

func init() {
	file.AddManager(StorageTypeFS, &FileManager{})
}
