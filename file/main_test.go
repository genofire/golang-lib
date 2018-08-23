package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTOML(t *testing.T) {
	assert := assert.New(t)

	a := struct {
		Text string `toml:"text"`
	}{}

	err := ReadTOML("testfiles/donoexists", &a)
	assert.Error(err, "could find file ^^")

	err = ReadTOML("testfiles/trash.txt", &a)
	assert.Error(err, "could marshel file ^^")

	err = ReadTOML("testfiles/ok.toml", &a)
	assert.NoError(err)
	assert.Equal("hallo", a.Text)
}

func TestReadJSON(t *testing.T) {
	assert := assert.New(t)

	a := struct {
		Text string `toml:"text"`
	}{}

	err := ReadJSON("testfiles/donoexists", &a)
	assert.Error(err, "could find file ^^")

	err = ReadJSON("testfiles/trash.txt", &a)
	assert.Error(err, "could marshel file ^^")

	err = ReadJSON("testfiles/ok.json", &a)
	assert.NoError(err)
	assert.Equal("hallo", a.Text)
}

func TestSaveJSON(t *testing.T) {
	assert := assert.New(t)

	tmpfile, _ := ioutil.TempFile("/tmp", "lib-json-testfile.json")
	err := SaveJSON(tmpfile.Name(), 3)
	assert.NoError(err, "could not save temp")

	err = SaveJSON(tmpfile.Name(), tmpfile.Name)
	assert.Error(err, "could not save func")

	err = SaveJSON("/proc/readonly", 4)
	assert.Error(err, "could not save to /dev/null")

	var testvalue int
	err = ReadJSON(tmpfile.Name(), &testvalue)
	assert.NoError(err)
	assert.Equal(3, testvalue)
	os.Remove(tmpfile.Name())
}
