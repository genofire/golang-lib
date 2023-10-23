package fs_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"codeberg.org/genofire/golang-lib/web/file"
	"codeberg.org/genofire/golang-lib/web/file/fs"
)

func TestOpenStat(t *testing.T) {
	assert := assert.New(t)
	var fs file.FS = &fs.FS{Root: "./testdata"}
	assert.NoError(fs.Check())

	f, err := fs.Open("glenda")
	assert.NoError(err)
	assert.NotNil(f)

	fi, err := f.Stat()
	assert.NoError(err)
	assert.NotNil(fi)

	assert.Equal(uuid.MustParse("d2750ced-4bdc-41d0-8c2f-5b9de44b84ef"), fi.(file.FileInfo).ID())
	assert.Equal("text/plain", fi.(file.FileInfo).ContentType())
	assert.Equal("glenda", fi.Name())
	assert.Equal(int64(99), fi.Size())
}

func TestCreateOpenUUIDRead(t *testing.T) {
	assert := assert.New(t)
	var fs file.FS = &fs.FS{Root: "./testdata"}
	assert.NoError(fs.Check())

	err := fs.Store(uuid.MustParse("f9375ccb-ee09-4ecf-917e-b88725efcb68"), "$name", "text/plain", strings.NewReader("hello, world\n"))
	assert.NoError(err)

	f, err := fs.OpenUUID(uuid.MustParse("f9375ccb-ee09-4ecf-917e-b88725efcb68"))
	assert.NoError(err)
	assert.NotNil(f)

	buf, err := io.ReadAll(f)
	assert.NoError(err)
	assert.Equal("hello, world\n", string(buf))

	os.RemoveAll("./testdata/f9375ccb-ee09-4ecf-917e-b88725efcb68")
}
