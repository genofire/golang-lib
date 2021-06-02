package database

import (
	"errors"
)

var (
	// ErrNotConnected - database is not connected
	ErrNotConnected = errors.New("database is not connected")
	// ErrNothingToMigrate if nothing has to be migrated
	ErrNothingToMigrate = errors.New("there is nothing to migrate")
)
