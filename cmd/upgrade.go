package main

import (
	"fmt"
	"log"
	"strings"

	"mercury/internal/migrations"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/lib/pq"
	"golang.org/x/mod/semver"
)

type migFunc struct {
	version string
	fn      func(*sqlx.DB, *koanf.Koanf, *log.Logger) error
}

var migList = []migFunc{
	{"v0.1.0", migrations.V0_1_0},
}

// upgrade upgrades the database to the current version by running SQL migration files
// for all version from the last known version to the current one.
func upgrade(db *sqlx.DB, prompt bool) {
	if prompt {
		var ok string
		fmt.Printf("** IMPORTANT: Take a backup of the database before upgrading.\n")
		fmt.Print("continue (y/n)?  ")
		if _, err := fmt.Scanf("%s", &ok); err != nil {
			log.Fatalf("error reading value from terminal: %v", err)
		}
		if strings.ToLower(ok) != "y" {
			fmt.Println("upgrade cancelled")
			return
		}
	}

	_, toRun, err := getPendingMigrations(db)
	if err != nil {
		log.Fatalf("error checking migrations: %v", err)
	}

	// No migrations to run.
	if len(toRun) == 0 {
		log.Printf("no upgrades to run. Database is up to date.")
		return
	}

	// Execute migrations in succession.
	for _, m := range toRun {
		log.Printf("running migration %s", m.version)
		if err := m.fn(db, ko, log); err != nil {
			log.Fatalf("error running migration %s: %v", m.version, err)
		}

		// Record the migration version in the settings table. There was no
		// settings table until v0.7.0, so ignore the no-table errors.
		if err := recordMigrationVersion(m.version, db); err != nil {
			if isTableNotExistErr(err) {
				continue
			}
			log.Fatalf("error recording migration version %s: %v", m.version, err)
		}
	}

	log.Printf("upgrade complete")
}

func checkUpgrade(db *sqlx.DB) {
	lastVer, toRun, err := getPendingMigrations(db)
	if err != nil {
		log.Fatalf("error checking migrations: %v", err)
	}

	// No migrations to run.
	if len(toRun) == 0 {
		return
	}

	var vers []string
	for _, m := range toRun {
		vers = append(vers, m.version)
	}

	log.Fatalf(`there are %d pending database upgrade(s): %v. The last upgrade was %s. Backup the database and run listmonk --upgrade`,
		len(toRun), vers, lastVer)
}

func getPendingMigrations(db *sqlx.DB) (string, []migFunc, error) {
	lastVer, err := getLastMigrationVersion(db)
	if err != nil {
		return "", nil, err
	}

	// Iterate through the migration versions and get everything above the last
	// upgraded semver.
	var toRun []migFunc
	for i, m := range migList {
		if semver.Compare(m.version, lastVer) > 0 {
			toRun = migList[i:]
			break
		}
	}

	return lastVer, toRun, nil
}

// getLastMigrationVersion returns the last migration semver recorded in the DB.
// If there isn't any, `v0.0.0` is returned.
func getLastMigrationVersion(db *sqlx.DB) (string, error) {
	var v string
	if err := db.Get(&v, `
		SELECT COALESCE(
			(SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1),
		'v0.0.0')`); err != nil {
		if isTableNotExistErr(err) {
			return "v0.0.0", nil
		}
		return v, err
	}
	return v, nil
}

// isTableNotExistErr checks if the given error represents a Postgres/pq
// "table does not exist" error.
func isTableNotExistErr(err error) bool {
	if p, ok := err.(*pq.Error); ok {
		// `settings` table does not exist. It was introduced in v0.7.0.
		if p.Code == "42P01" {
			return true
		}
	}
	return false
}
