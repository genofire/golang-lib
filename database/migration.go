package database

import (
	"sort"

	gormigrate "github.com/genofire/gormigrate/v2"
)

func (config *Database) sortedMigration(testdata bool) []*gormigrate.Migration {
	var migrations []*gormigrate.Migration
	for _, v := range config.migrations {
		migrations = append(migrations, v)
	}
	if testdata {
		for _, v := range config.migrationTestdata {
			migrations = append(migrations, v)
		}
	}
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})
	return migrations
}

func (config *Database) setupMigrator(testdata bool) (*gormigrate.Gormigrate, error) {
	migrations := config.sortedMigration(testdata)
	if len(migrations) == 0 {
		return nil, ErrNothingToMigrate
	}

	return gormigrate.New(config.DB, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            true,
		ValidateUnknownMigrations: false,
	}, migrations), nil

}

func (config *Database) migrate(testdata bool) error {
	m, err := config.setupMigrator(testdata)
	if err != nil {
		return err
	}
	return m.Migrate()
}

// Migrate run migration
func (config *Database) Migrate() error {
	return config.migrate(false)
}

// MigrateTestdata run migration and testdata migration
func (config *Database) MigrateTestdata() error {
	return config.migrate(true)
}

// AddMigration add to database config migration step
func (config *Database) AddMigration(m ...*gormigrate.Migration) {
	config.addMigrate(false, m...)
}

// AddMigrationTestdata add to database config migration step of testdata
func (config *Database) AddMigrationTestdata(m ...*gormigrate.Migration) {
	config.addMigrate(true, m...)
}
func (config *Database) addMigrate(testdata bool, m ...*gormigrate.Migration) {
	if config.migrations == nil {
		config.migrations = make(map[string]*gormigrate.Migration)
	}
	if config.migrationTestdata == nil {
		config.migrationTestdata = make(map[string]*gormigrate.Migration)
	}

	for _, i := range m {
		if testdata {
			config.migrationTestdata[i.ID] = i
		} else {
			config.migrations[i.ID] = i
		}
	}
}

// ReMigrate Rollback und run every migration step again till id
func (config *Database) ReMigrate(id string) error {
	migrations := config.sortedMigration(true)
	m, err := config.setupMigrator(true)
	if err != nil {
		return err
	}

	x := 0
	for _, m := range migrations {
		if m.ID == id {
			break
		}
		x = x + 1
	}
	// TODO not found

	for i := len(migrations) - 1; i >= x; i = i - 1 {
		mStep := migrations[i]
		if err := m.RollbackTo(mStep.ID); err != nil {
			return err
		}
	}
	return m.Migrate()
}
