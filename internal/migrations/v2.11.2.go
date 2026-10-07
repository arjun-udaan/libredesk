package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V2_11_2 adds managed attributes to desks that already applied v2.9.
func V2_11_2(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	_, err := db.Exec(`ALTER TABLE custom_attribute_definitions ADD COLUMN IF NOT EXISTS read_only BOOLEAN NOT NULL DEFAULT false;`)
	return err
}
