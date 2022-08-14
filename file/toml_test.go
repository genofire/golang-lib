package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOMLDuration(t *testing.T) {
	assert := assert.New(t)

	var d TOMLDuration
	err := d.UnmarshalText([]byte("5m"))
	assert.NoError(err)

	err = d.UnmarshalText([]byte("5z"))
	assert.Error(err)

	txt, err := d.MarshalText()
	assert.NoError(err)
	assert.Equal("5m0s", string(txt))
}

func TestReadTOML(t *testing.T) {
	assert := assert.New(t)

	a := struct {
		Text string `config:"text" toml:"text"`
	}{}

	err := ReadTOML("testfiles/donoexists", &a)
	assert.Error(err, "could find file ^^")

	err = ReadTOML("testfiles/trash.txt", &a)
	assert.Error(err, "could marshel file ^^")

	err = ReadTOML("testfiles/ok.toml", &a)
	assert.NoError(err)
	assert.Equal("hallo", a.Text)
}

func TestSaveTOML(t *testing.T) {
	assert := assert.New(t)

	type to struct {
		Value int `config:"v" toml:"v"`
	}
	toSave := to{Value: 3}

	tmpfile, _ := ioutil.TempFile("/tmp", "lib-json-testfile.json")
	err := SaveTOML(tmpfile.Name(), &toSave)
	assert.NoError(err, "could not save temp")

	err = SaveTOML(tmpfile.Name(), 3)
	assert.Error(err, "could not save func")

	toSave.Value = 4
	err = SaveTOML("/proc/readonly", &toSave)
	assert.Error(err, "could not save to /dev/null")

	testvalue := to{}
	err = ReadTOML(tmpfile.Name(), &testvalue)
	assert.NoError(err)
	assert.Equal(3, testvalue.Value)
	os.Remove(tmpfile.Name())
}
