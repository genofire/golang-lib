package file_test

import (
	"testing"

	fsfile "dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/web/file"
	"github.com/stretchr/testify/assert"
)

func TestCreateFSOK(t *testing.T) {
	assert := assert.New(t)

	config := file.FSInfo{}
	err := fsfile.ReadTOML("testdata/createfs_fs.toml", &config)
	assert.NoError(err)

	fs, err := config.Create()
	assert.NoError(err)
	assert.NoError(fs.Check())
}

func TestCreateS3(t *testing.T) {
	assert := assert.New(t)

	config := file.FSInfo{}
	err := fsfile.ReadTOML("testdata/createfs_s3.toml", &config)
	assert.NoError(err)

	fs, err := config.Create()
	assert.NoError(err)
	assert.NoError(fs.Check())
}

func TestCreateFSNotOK(t *testing.T) {
	assert := assert.New(t)

	config := file.FSInfo{}
	err := fsfile.ReadTOML("testdata/createfs_fsnone.toml", &config)
	assert.NoError(err)

	_, err = config.Create()
	assert.ErrorIs(err, file.ErrNoFSRoot)
}

func TestCreateFSNone(t *testing.T) {
	assert := assert.New(t)

	config := file.FSInfo{}
	err := fsfile.ReadTOML("testdata/createfs_none.toml", &config)

	// https://github.com/naoina/toml/pull/51
	assert.Contains(err.Error(), file.ErrInvalidFSType.Error())
}

func TestCreateFSInvalid(t *testing.T) {
	assert := assert.New(t)

	config := file.FSInfo{}
	_, err := config.Create()
	assert.ErrorIs(err, file.ErrNoFSRoot)
}
