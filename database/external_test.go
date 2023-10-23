package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"codeberg.org/genofire/golang-lib/database"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	// DBConnection = "user=root password=root dbname=defaultdb host=localhost port=26257 sslmode=disable"
	DBConnection = "user=root password=root dbname=defaultdb host=localhost sslmode=disable"
)

func ParseConfig(conn string) *database.Database {
	k := koanf.New("/")

	out := database.Database{}

	k.Load(rawbytes.Provider([]byte(conn)), toml.Parser())

	k.UnmarshalWithConf("", &out, koanf.UnmarshalConf{Tag: "config"})
	return &out
}

func TestConn(t *testing.T) {
	assert := assert.New(t)

	d := ParseConfig(`
	[connection]
	hostname = "localhost"
	username = "root"
	password = "a"
	dbname = "database"
	extra_options = "sslmode=disable"
	`)

	assert.Equal("postgresql://root:a@localhost/database?sslmode=disable", d.Connection.String(), "splitted")

	d = ParseConfig(`
	[connection]
	string = "postgresql://root:a@localhost/database?sslmode=disable"
	`)
	assert.Equal("postgresql://root:a@localhost/database?sslmode=disable", d.Connection.String(), "connection_string")

	d = ParseConfig(`
	[connection]
	string = "postgresql://root:a@localhost/database?sslmode=disable"
	username = "user"
	password = "b"
	`)
	assert.Equal("postgresql://user:b@localhost/database?sslmode=disable", d.Connection.String(), "both")

}
