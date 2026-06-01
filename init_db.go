package main

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

func initDB(ctx context.Context)(*sql.DB){
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

	// schema := 
	// db.ExecContext(ctx, schema)

	return db
}