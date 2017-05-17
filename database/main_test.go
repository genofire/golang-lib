// Package that provides the functionality to open, close and use a database
package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID    int64
	Value string `gorm:"type:varchar(255);column:value" json:"value"`
}

// Function to test the error handling for the database opening
func TestOpenNoDB(t *testing.T) {
	assert := assert.New(t)

	c := Config{}

	err := Open(c)
	assert.Error(err, "error")
}

// Function to test the opening of one database
func TestOpenOneDB(t *testing.T) {
	assert := assert.New(t)
	AddModel(&TestModel{})

	c := Config{
		Type:       "sqlite3",
		Logging:    true,
		Connection: "file:database?mode=memory",
	}
	var count int64

	err := Open(c)
	assert.NoError(err, "no error")

	Write.Create(&TestModel{Value: "first"})
	Write.Create(&TestModel{Value: "secound"})

	var list []*TestModel
	Read.Find(&list).Count(&count)
	assert.Equal(int64(2), count, "not enought entries")
	Close()
}

// Function to test the opening of a second database
func TestOpenTwoDB(t *testing.T) {
	assert := assert.New(t)
	AddModel(&TestModel{})
	c := Config{
		Type:           "sqlite3",
		Logging:        true,
		Connection:     "file:database?mode=memory",
		ReadConnection: "file/",
	}

	err := Open(c)
	assert.Error(err, "no error found")

	c = Config{
		Type:           "sqlite3",
		Logging:        true,
		Connection:     "file:database?mode=memory",
		ReadConnection: "file:database2?mode=memory",
	}
	var count int64

	err = Open(c)
	assert.NoError(err, "no error")

	Write.Create(&TestModel{Value: "first"})
	Write.Create(&TestModel{Value: "secound"})

	var list []*TestModel
	Write.Find(&list).Count(&count)
	assert.Equal(int64(2), count, "not enought entries")

	result := Read.Find(&list)
	assert.Error(result.Error, "error, because it is the wrong database")
	Close()
}
