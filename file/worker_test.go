package file

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveJSONWorker(t *testing.T) {
	assert := assert.New(t)

	tmpfile, _ := ioutil.TempFile("/tmp", "lib-json-workertest.json")

	worker := NewSaveJSONWorker(100*time.Millisecond, tmpfile.Name(), 12)
	assert.NotNil(worker)

	time.Sleep(300 * time.Millisecond)

	var testvalue int
	err := ReadJSON(tmpfile.Name(), &testvalue)
	assert.NoError(err)
	assert.Equal(12, testvalue)
	os.Remove(tmpfile.Name())
}
