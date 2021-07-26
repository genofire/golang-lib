package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/web/file2"
)

// TODO: actually test, either using little dummies or using sth like play.min.io

func TestTypes(t *testing.T) {
	assert := assert.New(t)

	var fs file.FS
	fs, err := New("127.0.0.1", false, "", "", "", "")
	_ = fs
	assert.Error(err)
}

func ExampleNew() {
	New("play.min.io", true, "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "file-store", "")
}
