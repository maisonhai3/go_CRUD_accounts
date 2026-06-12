package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"accountCRUD/handlers"
	"accountCRUD/repositories"
)

func main() {
	db := repositories.InitDB(context.Background())
	repo := &repositories.DBHandler{
		DBConn: db,
	}
	defer db.Close()

	h := &handlers.Handler{Repo: repo}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts", h.GetAccounts)
	mux.HandleFunc("GET /accounts/{id}", h.GetAccountById)
	mux.HandleFunc("POSt /accounts", h.CreateAccount)

	// Config this server Manually
	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown

	// Just monitoring if we are running properly
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Waiting for OS say us to cook
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // We're waiting

	// Okay, OS said it. I quit.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown Failed: %v", err)
	}
}
