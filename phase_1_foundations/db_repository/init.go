package dbrepository

import (
	"context"
	"database/sql"
	"log"
)

func Init(ctx context.Context) DBRepo {
	db, err := sql.Open("sqlite", "file:accounts.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("Failed to init DB: ", err.Error())
	}

	// Verify
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to init DB: ", err.Error())
	}

	// Migrate
	Migrate(ctx, db)

	// Seed

	return DBRepo{
		Conn: db,
	}
}
