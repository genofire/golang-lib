package database

import (
	gormigrate "github.com/genofire/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database struct to read from config
type Database struct {
	DB                *gorm.DB        `toml:"-"`
	Connection        string          `toml:"connection"`
	Debug             bool            `toml:"debug"`
	Testdata          bool            `toml:"testdata"`
	LogLevel          logger.LogLevel `toml:"log_level"`
	migrations        map[string]*gormigrate.Migration
	migrationTestdata map[string]*gormigrate.Migration
}

// Run database config - connect and migrate
func (config *Database) Run() error {
	if err := config.run(); err != nil {
		return err
	}
	return config.migrate(config.Testdata)
}

// ReRun database config - connect and  re migration
func (config *Database) ReRun() error {
	if err := config.run(); err != nil {
		return err
	}
	m, err := config.setupMigrator(true)
	if err != nil {
		return err
	}
	if err := m.RollbackAll(); err != nil {
		return err
	}
	return config.migrate(config.Testdata)
}

func (config *Database) run() error {
	db, err := gorm.Open(postgres.Open(config.Connection), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return err
	}
	db.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if config.Debug {
		db = db.Debug()
	}

	config.DB = db
	return nil
}

// Status get status - is database pingable
func (config *Database) Status() error {
	if config.DB == nil {
		return ErrNotConnected
	}
	sqlDB, err := config.DB.DB()
	if err != nil {
		return err

	}
	if err = sqlDB.Ping(); err != nil {
		return err
	}
	return nil
}
