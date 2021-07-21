package fs

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/file"
)

func TestCheck(t *testing.T) {
	assert := assert.New(t)

	service := file.Service{
		StorageType: StorageTypeFS,
		Path:        "./test",
	}

	assert.NoError(service.Check())

	service.StorageType = "s3"
	assert.ErrorIs(file.ErrUnsupportedStorageType, service.Check())

	service.StorageType = StorageTypeFS
	service.Path = "./main_test.go"
	assert.ErrorIs(ErrPathNotExistsOrNoDirectory, service.Check())

	/* TODO no write permission
	service.Path = "/dev"
	assert.ErrorIs(ErrPathNotExistsOrNoDirectory, service.Check())
	*/
}

func TestSave(t *testing.T) {
	assert := assert.New(t)

	service := file.Service{
		StorageType: "s3",
		Path:        "./test",
	}

	_, err := service.Upload(nil)
	assert.ErrorIs(file.ErrUnsupportedStorageType, err)

	service.StorageType = StorageTypeFS
	req, err := web.NewRequestWithFile("localhost", "./test/00000000-0000-0000-0000-000000000000/a.txt")
	assert.NoError(err)
	assert.NotNil(req)

	_, err = service.Upload(req)
	assert.NoError(err)

	service.Path = "/dev"
	_, err = service.Upload(req)
	assert.True(os.IsNotExist(err))
	//assert.True(os.IsPermission(err))

	// TODO no write permission
}

func TestRead(t *testing.T) {
	assert := assert.New(t)

	service := file.Service{
		StorageType: "s3",
		Path:        "./test",
	}

	_, err := service.Read(nil)
	assert.ErrorIs(file.ErrUnsupportedStorageType, err)

	service.StorageType = StorageTypeFS

	file := &file.File{
		Path: "00000000-0000-0000-0000-000000000000/a.txt",
	}
	r, err := service.Read(file)
	assert.NoError(err)
	buf := &strings.Builder{}
	_, err = io.Copy(buf, r)
	assert.Equal("Hello world\n", buf.String())

	service.Path = "/dev"
	_, err = service.Read(file)
	assert.True(os.IsNotExist(err))

	// TODO no write permission
}
