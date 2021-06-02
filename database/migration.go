package database

import (
	"sort"

	"github.com/bdlm/log"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
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

// Migrate run migration
func (config *Database) Migrate() error {
	return config.migrate(false)
}

// MigrateTestdata run migration and testdata migration
func (config *Database) MigrateTestdata() error {
	return config.migrate(true)
}
func (config *Database) migrate(testdata bool) error {
	migrations := config.sortedMigration(testdata)
	if len(migrations) == 0 {
		return ErrNothingToMigrate
	}

	m := gormigrate.New(config.DB, gormigrate.DefaultOptions, migrations)

	return m.Migrate()
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
			config.migrations[i.ID] = i
		} else {
			config.migrationTestdata[i.ID] = i
		}
	}
}

// ReRun Rollback und run every migration step again till id
func (config *Database) ReRun(id string) {
	migrations := config.sortedMigration(true)
	x := 0
	for _, m := range migrations {
		if m.ID == id {
			break
		}
		x = x + 1
	}

	for i := len(migrations) - 1; i >= x; i = i - 1 {
		m := migrations[i]
		log.Warnf("rollback %s", m.ID)
		err := m.Rollback(config.DB)
		if err != nil {
			log.Errorf("rollback %s", err)
		}
	}
	for _, m := range migrations {
		log.Warnf("run %s", m.ID)
		err := m.Migrate(config.DB)
		if err != nil {
			log.Errorf("run %s", err)
		}
	}
}
