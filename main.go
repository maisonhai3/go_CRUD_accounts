package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func getAccountById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	fmt.Printf("ctx channel: %v\n", ctx.Done())

	type ctxKey string
	const userIDKey ctxKey = "user-id"
	id := r.PathValue("id")
	ctx = context.WithValue(ctx, userIDKey, id)

	// Fake a DB call
	func(ctx context.Context, id string) {
		fmt.Printf("ctx channel: %v\n", ctx.Done())
		log.Printf("Getting Account By Id %v", id)

		select {
		case <-time.After(2 * time.Second):
			log.Printf("Finished fetching %s", id)
		case <-ctx.Done():
			log.Printf("Abort fetching %s: %v", id, ctx.Err())
		}
	}(ctx, id)

	if err := ctx.Err(); err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	for _, a := range accounts {
		if a.ID == id {
			aJSON, err := json.Marshal(a)
			if err != nil {
				http.Error(w, "Error in marshalling", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(aJSON)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
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
	mux.HandleFunc("GET /accounts/{id}", getAccountById)

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
