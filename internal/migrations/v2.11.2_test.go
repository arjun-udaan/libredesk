package migrations

import (
	"github.com/abhinavxd/libredesk/internal/testutil"
	"testing"
)

func TestV2_11_2ExistingDesk(t *testing.T) {
	db := testutil.NewDB(t, "migration_v2_11_2")
	db.MustExec(`ALTER TABLE custom_attribute_definitions DROP COLUMN read_only`)
	db.MustExec(`INSERT INTO custom_attribute_definitions (name, description, applies_to, key, data_type) VALUES ('Tier', 'Account tier', 'contact', 'migration_tier', 'text')`)
	for range 2 {
		if err := V2_11_2(db, nil, nil); err != nil {
			t.Fatal(err)
		}
	}
	var readOnly bool
	if err := db.Get(&readOnly, `SELECT read_only FROM custom_attribute_definitions WHERE key='migration_tier'`); err != nil || readOnly {
		t.Fatalf("read_only=%v, error=%v", readOnly, err)
	}
	db.MustExec(`UPDATE custom_attribute_definitions SET read_only=true WHERE key='migration_tier'`)
	if err := V2_11_2(db, nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := db.Get(&readOnly, `SELECT read_only FROM custom_attribute_definitions WHERE key='migration_tier'`); err != nil || !readOnly {
		t.Fatalf("saved read_only=%v, error=%v", readOnly, err)
	}
}
