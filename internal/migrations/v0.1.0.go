package migrations

import (
	"log"
	"mercury/internal/config"

	"github.com/jmoiron/sqlx"
)

func V0_1_0(db *sqlx.DB, config *config.Config, log *log.Logger) error {
	log.Print("Running migration v0.1.0")
	return nil
}
