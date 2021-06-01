// Package database provides the functionality to open, close and use a database
package database

import (
	"gorm.io/gorm"
	// load gorm defaults driver
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"

	"github.com/bdlm/log"
)

// Write Database connection for writing purposes
var Write *gorm.DB

// Read Database connection for reading purposes
var Read *gorm.DB

// Configuration files
var (
	config  *Config
	runtime []interface{}
)

// Config of the database connection
type Config struct {
	// type of the database, currently supports sqlite and postgres
	Type string
	// connection configuration
	Connection string
	// create another connection for reading only
	ReadConnection string
	// enable logging of the generated sql string
	Logging bool
}

// Open database and set the given configuration
func Open(c Config) (err error) {
	writeLog := log.WithField("db", "write")
	config = &c
	Write, err = gorm.Open(config.Type, config.Connection)
	if err != nil {
		return
	}
	Write.SingularTable(true)
	Write.LogMode(c.Logging)
	Write.SetLogger(writeLog)
	Write.Callback().Create().Remove("gorm:update_time_stamp")
	Write.Callback().Update().Remove("gorm:update_time_stamp")
	if len(config.ReadConnection) > 0 {
		readLog := log.WithField("db", "read")
		Read, err = gorm.Open(config.Type, config.ReadConnection)
		if err != nil {
			return
		}
		Read.SingularTable(true)
		Read.LogMode(c.Logging)
		Read.SetLogger(readLog)
		Read.Callback().Create().Remove("gorm:update_time_stamp")
		Read.Callback().Update().Remove("gorm:update_time_stamp")
	} else {
		Read = Write
	}
	Write.AutoMigrate(runtime...)
	return
}

// Close connnection to database safely
func Close() {
	Write.Close()
	if len(config.ReadConnection) > 0 {
		Read.Close()
	}
}

// AddModel to the runtime
func AddModel(m interface{}) {
	runtime = append(runtime, m)
}
