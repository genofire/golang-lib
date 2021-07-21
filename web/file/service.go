package file

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// A Service to handle file-uploads in golang
type Service struct {
	StorageType string `toml:"storage_type"`
	Path        string `toml:"path"`
}

// Check if Service is configurated and useable
func (s *Service) Check() error {
	mgmt, ok := managers[s.StorageType]
	if !ok {
		return ErrUnsupportedStorageType
	}
	return mgmt.Check(s)
}

// Upload a file to storage
func (s *Service) Upload(request *http.Request) (*File, error) {
	mgmt, ok := managers[s.StorageType]
	if !ok {
		return nil, ErrUnsupportedStorageType
	}
	file, fileRequest, err := request.FormFile("file")
	if err != nil {
		return nil, err
	}
	fileObj := File{
		Filename: filepath.Base(fileRequest.Filename),
	}

	// detect contenttype
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	fileObj.ContentType = http.DetectContentType(buffer[:n])

	// Reset the read pointer
	file.Seek(0, io.SeekStart)
	if err := mgmt.Save(s, &fileObj, file); err != nil {
		return nil, err
	}
	return &fileObj, nil

}

// GINUpload a file to storage using gin-gonic
func (s *Service) GINUpload(c *gin.Context) (*File, error) {
	return s.Upload(c.Request)
}

// Read a file to storage
func (s *Service) Read(file *File) (io.Reader, error) {
	mgmt, ok := managers[s.StorageType]
	if !ok {
		return nil, ErrUnsupportedStorageType
	}
	return mgmt.Read(s, file)
}
