package migrations

import (
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V2_11_1 applies merged PR changes to desks that already ran the v2.9 migration.
func V2_11_1(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf) error {
	if _, err := db.Exec(`
		INSERT INTO settings ("key", value) VALUES
			('app.reply_guard_phrases', '""'::jsonb),
			('app.time_format', '"12h"'::jsonb)
		ON CONFLICT ("key") DO NOTHING;
	`); err != nil {
		return err
	}
	return migrateExternalSync(db)
}
