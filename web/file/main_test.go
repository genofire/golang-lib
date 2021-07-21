package file

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"dev.sum7.eu/genofire/golang-lib/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const storageTypeDummy = "dummy"

type dummyManager struct {
}

func (m *dummyManager) Check(s *Service) error {
	return nil
}
func (m *dummyManager) Save(s *Service, file *File, src multipart.File) error {
	if src == nil {
		return errors.New("nothing to fill")
	}
	return nil
}
func (m *dummyManager) Read(s *Service, file *File) (io.Reader, error) {
	b := bytes.Buffer{}
	b.WriteString("Hello world\n")
	return &b, nil
}

func init() {
	AddManager(storageTypeDummy, &dummyManager{})
}

func TestCheck(t *testing.T) {
	assert := assert.New(t)

	service := Service{
		StorageType: storageTypeDummy,
		Path:        "./fs/test",
	}
	assert.NoError(service.Check())

	service.StorageType = "s3"
	assert.ErrorIs(ErrUnsupportedStorageType, service.Check())
}

func TestSave(t *testing.T) {
	assert := assert.New(t)

	service := Service{
		StorageType: "fs",
		Path:        "./fs/test",
	}

	_, err := service.Upload(nil)
	assert.ErrorIs(ErrUnsupportedStorageType, err)

	service.StorageType = storageTypeDummy
	_, err = service.GINUpload(&gin.Context{Request: &http.Request{}})
	assert.ErrorIs(err, http.ErrNotMultipart)

	req, err := web.NewRequestWithFile("http://localhost/upload", "./fs/test/00000000-0000-0000-0000-000000000000/a.txt")
	assert.NoError(err)
	assert.NotNil(req)

	_, err = service.Upload(req)
	assert.NoError(err)
}

func TestRead(t *testing.T) {
	assert := assert.New(t)

	service := Service{
		StorageType: "fs",
		Path:        "./fs/test",
	}

	_, err := service.Read(nil)
	assert.ErrorIs(ErrUnsupportedStorageType, err)

	service.StorageType = "dummy"

	file := &File{
		Path: "00000000-0000-0000-0000-000000000000/a.txt",
	}
	r, err := service.Read(file)
	assert.NoError(err)
	buf := &strings.Builder{}
	_, err = io.Copy(buf, r)
	assert.Equal("Hello world\n", buf.String())
}
