package migrations

import (
	"testing"

	"github.com/abhinavxd/libredesk/internal/testutil"
)

func TestV2_11_1ExistingDesk(t *testing.T) {
	db := testutil.NewDB(t, "migration_v2_11_1")
	db.MustExec(`ALTER TABLE users DROP COLUMN external_sync`)
	db.MustExec(`DELETE FROM settings WHERE key IN ('app.reply_guard_phrases', 'app.time_format')`)
	db.MustExec(`INSERT INTO users (type, email, first_name, external_user_id) VALUES ('contact', 'sync@example.test', 'Synced', 'external-test')`)
	for range 2 {
		if err := V2_11_1(db, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
	var count int
	if err := db.Get(&count, `SELECT count(*) FROM settings WHERE key IN ('app.reply_guard_phrases', 'app.time_format')`); err != nil || count != 2 {
		t.Fatalf("settings count = %d, error = %v", count, err)
	}
	var name string
	if err := db.Get(&name, `SELECT external_sync->>'first_name' FROM users WHERE external_user_id = 'external-test'`); err != nil || name != "Synced" {
		t.Fatalf("synced name = %q, error = %v", name, err)
	}
	db.MustExec(`UPDATE settings SET value = '"24h"'::jsonb WHERE key = 'app.time_format'`)
	if err := V2_11_1(db, nil, nil); err != nil {
		t.Fatal(err)
	}
	var value string
	if err := db.Get(&value, `SELECT value #>> '{}' FROM settings WHERE key = 'app.time_format'`); err != nil || value != "24h" {
		t.Fatalf("saved time format = %q, error = %v", value, err)
	}
}
