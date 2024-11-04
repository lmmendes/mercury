package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func install(db *sqlx.DB, prompt, idempotent bool) {
	// Check if the database is already initialized.
	// If the database is not initialized, we should get "v0.0.0" as the version.
	version, err := getLastMigrationVersion(db)
	if err != nil {
		log.Fatalf("error getting last migration version: %v", err)
	}

	if version != "v0.0.0" {
		log.Fatalf("database is already initialized. Current version is %s", version)
	}

	// Fetch all available migrations and run them.
	lastVer, toRun, err := getPendingMigrations(db)
	if err != nil {
		log.Fatalf("error checking migrations: %v", err)
	}

	// No migrations to run.
	if len(toRun) == 0 {
		return
	} else {
	}

	var vers []string
	for _, m := range toRun {
		vers = append(vers, m.version)
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

}

func checkSchema(db *sqlx.DB) (bool, error) {
	if _, err := db.Exec(`SELECT version FROM schema_migrations LIMIT 1`); err != nil {
		if isTableNotExistErr(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
