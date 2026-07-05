package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	dbrepository "phase_1_foundations/db_repository"
	httphandler "phase_1_foundations/http_handler"
	"syscall"
	"time"
)

func main() {
	// Init DBMS
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbRepo := dbrepository.Init(ctx)
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dbRepo.Conn)

	// Init Server
	// Routing
	httpHandler := httphandler.HTTPHandler{DBRepo: dbRepo}
	muxServer := http.NewServeMux()
	muxServer.HandleFunc("POST /accounts", httpHandler.CreateAccount)
	// Config
	httpSrv := http.Server{
		Addr:         ":8080",
		Handler:      muxServer,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Gracefully shutdown
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown gracefully: %v", err.Error())
	}
}
