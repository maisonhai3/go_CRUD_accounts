package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type DBHandler struct {
	DBConn *sql.DB
}

func InitDB(ctx context.Context) *sql.DB {
	// db, err := sql.Open("sqlite", "file:accounts.db")
	db, err := sql.Open("sqlite", "file:accounts.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("DB init failed: %v", err.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("DB init failed: %v", err.Error())
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS accounts (
               id         TEXT PRIMARY KEY,
            name       TEXT NOT NULL,
            currency   TEXT NOT NULL,        -- ISO 4217: "USD", "VND", "JPY"
            balance    INTEGER NOT NULL,     -- minor units. KHÔNG BAO GIỜ REAL/float.
            created_at TEXT NOT NULL,        -- RFC3339
            updated_at TEXT NOT NULL,
            deleted_at TEXT                  -- NULL = chưa xóa (soft delete)
    
)`)

	if err != nil {
		log.Fatalf("DB init failed: %v", err.Error())
	}

	return db
}
