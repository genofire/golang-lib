package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	// DBConnection = "postgresql://root:root@localhost:26257/defaultdb?sslmode=disable"
	DBConnection = "postgresql://root:root@localhost/defaultdb?sslmode=disable"
)

func TestStatus(t *testing.T) {
	assert := assert.New(t)

	d := Database{
		Debug: true,
	}
	d.Connection.URI = "postgresql://localhost"
	err := d.Status()
	assert.Error(err)
	assert.Equal(ErrNotConnected, err)

	err = d.Run()
	assert.Error(err)
	assert.Contains(err.Error(), "failed to connect")

	d.Connection.URI = DBConnection
	err = d.Run()
	assert.Error(err)
	assert.Equal(ErrNothingToMigrate, err)

	err = d.Status()
	assert.NoError(err)
}
