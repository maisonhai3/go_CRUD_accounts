package dbrepository

import (
	"context"
	"database/sql"
	"log"
)

func Migrate(ctx context.Context, db *sql.DB) {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS accounts (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            name       TEXT NOT NULL,
            currency   TEXT NOT NULL,        -- ISO 4217: "USD", "VND", "JPY"
            balance    INTEGER NOT NULL,     -- minor units. KHÔNG BAO GIỜ REAL/float.
            created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),        -- RFC3339
            updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ','now')),
            deleted_at TEXT                  -- NULL = chưa xóa (soft delete)
	)`)

	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
