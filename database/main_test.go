package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	DBConnection = "user=root password=root dbname=defaultdb host=localhost port=26257 sslmode=disable"
)

func TestStatus(t *testing.T) {
	assert := assert.New(t)

	d := Database{
		Debug: true,
	}
	err := d.Status()
	assert.Error(err)
	assert.Equal(ErrNotConnected, err)

	err = d.Run()
	assert.Error(err)
	assert.Contains(err.Error(), "dial error")

	d.Connection = DBConnection
	err = d.Run()
	assert.Error(err)
	assert.Equal(ErrNothingToMigrate, err)

	err = d.Status()
	assert.NoError(err)
}
