package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type account struct {
	ID        string
	Name      string
	Currency  string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

var accounts = []account{
	{
		ID:        "123",
		Name:      "Hoang Hai Ha Van",
		Currency:  "USD",
		Balance:   1000,
		CreatedAt: time.Now(),
	},
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	buf, err := json.Marshal(accounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts", getAccounts)

	// Config this server Manually
	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // We're waiting

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown Failed: %v", err)
	}
}
