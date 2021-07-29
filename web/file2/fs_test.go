package file_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web/file2"
)

type TestFS struct {
	assert      *assert.Assertions
	filename    string
	data        string
	contentType string
}

func (f TestFS) Store(id uuid.UUID, name, contentType string, data io.Reader) error {
	f.assert.Equal(f.filename, name)
	dat, err := io.ReadAll(data)
	f.assert.NoError(err)
	f.assert.Equal(f.data, string(dat))
	contentType, _, err = mime.ParseMediaType(contentType)
	f.assert.NoError(err)
	f.assert.Equal(f.contentType, contentType)
	return nil
}

func (f TestFS) RemoveUUID(id uuid.UUID) error {
	return errors.New("TestFS.RemoveUUID called")
}

func (f TestFS) Open(name string) (fs.File, error) {
	return nil, errors.New("TestFS.Open called")
}

func (f TestFS) OpenUUID(uuid.UUID) (fs.File, error) {
	return nil, errors.New("TestFS.OpenUUID called")
}

func (f TestFS) Check() error { return nil }

func TestStoreFromHTTP(t *testing.T) {
	assert := assert.New(t)

	testfs := TestFS{
		assert:      assert,
		filename:    "cute-cat.png",
		data:        "content\nof file",
		contentType: "text/plain",
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	rq := httptest.NewRequest("PUT", "/", r)
	rq.Header.Set("Content-Type", m.FormDataContentType())
	go func() {
		f, err := m.CreateFormFile("file", testfs.filename)
		assert.NoError(err)
		_, err = f.Write([]byte(testfs.data))
		assert.NoError(err)
		m.Close()
	}()

	assert.NoError(file.StoreFromHTTP(testfs, rq, "file"))
}

var fstore file.FS

func ExampleFS() {
	// generate the UUID for the new file
	id := uuid.New()

	// store a file
	{
		f, _ := os.Open("glenda.png")
		fstore.Store(id, "glenda.png", "image/png", f)
		f.Close()
	}

	// copy back to a local file
	{
		r, _ := fstore.OpenUUID(id)
		w, _ := os.Create("glenda.png")
		io.Copy(w, r)
		r.Close()
		w.Close()
	}
}

func ExampleStoreFromHTTP() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if err := file.StoreFromHTTP(fstore, r, "file"); err != nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message":"%v"}`, err)))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}
