package dbrepository

import (
	"context"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Init(ctx context.Context) DBRepo {
	db, err := sql.Open("sqlite", "file:accounts.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	// Verify
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	// Migrate
	Migrate(ctx, db)

	// Seed

	return DBRepo{
		Conn: db,
	}
}
