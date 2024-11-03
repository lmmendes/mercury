package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
)

func V0_1_0(db *sqlx.DB, ko *koanf.Koanf, log *log.Logger) error {
	log.Print("Running migration v0.1.0")
	return nil
}
