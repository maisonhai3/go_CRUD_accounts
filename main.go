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

func getAccountById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	ctx = context.WithValue(ctx, "userId", id)

	// Fake a DB call
	go func(ctx context.Context) {
		log.Printf("Get Account By Id %v", ctx.Value("userId"))
		time.Sleep(1 * time.Second)
	}(ctx)

	for _, a := range accounts {
		if a.ID == id {
			aJSON, err := json.Marshal(a)
			if err != nil {
				log.Fatalf("Error in marshalling: %v", err.Error())
				return
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(aJSON)
			return
		}
	}
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
	mux.HandleFunc("GET /accounts/:id", getAccounts)

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
