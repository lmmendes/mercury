package main

import (
	"fmt"
	"log"
	"mercury/internal/config"

	"github.com/jmoiron/sqlx"
)

func install(db *sqlx.DB, config *config.Config, prompt, idempotent bool) {
	// Check if the database is already initialized.
	// If the database is not initialized, we should get "v0.0.0" as the version.
	version, err := getLastMigrationVersion(db)
	if err != nil {
		logger.Fatalf("Error getting last migration version: %v", err)
	}

	if version != "v0.0.0" {
		logger.Fatalf("Database is already initialized. Current version is %s", version)
	}

	// Fetch all available migrations and run them.
	_, toRun, err := getPendingMigrations(db)
	if err != nil {
		logger.Fatalf("Error checking migrations: %v", err)
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

	for _, m := range toRun {
		log.Printf("Running migration %s", m.version)
		if err := m.fn(db, config, logger); err != nil {
			log.Fatalf("Error running migration %s: %v", m.version, err)
		}

		if err := recordMigrationVersion(m.version, db); err != nil {
			log.Fatalf("Error recording migration version %s: %v", m.version, err)
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

func checkInstall(db *sqlx.DB) {
	if ok, err := checkSchema(db); err != nil {
		logger.Fatalf("error checking schema in DB: %v", err)
	} else if !ok {
		logger.Fatal("The database does not appear to be setup. Run --install.")
	}
}

func recordMigrationVersion(version string, db *sqlx.DB) error {
	_, err := db.Exec(fmt.Sprintf(`INSERT INTO schema_migrations (version) VALUES('%s')`, version))
	return err
}
