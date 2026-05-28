package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type account struct {
	ID         string
	Name       string
	Currency   string
	Balance    int64
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

var accounts = []account{
	{
		ID:         "123",
		Name:       "Hoang Hai Ha Van",
		Currency:   "USD",
		Balance:    1000,
		Created_at: time.Now(),
	},
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r) // prints the pointer address
	fmt.Println("Method:", r.Method)
	fmt.Println("URL:", r.URL)
	fmt.Println("Headers:", r.Header)

	fmt.Fprintf(w, "Write") // Use json.Marshal()
}

func main() {
	http.HandleFunc("/accounts", getAccounts)
	 
	// Config this server Manually
	srv := http.Server{
		Addr: 8080,
		Handler: mux,
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}
	log.Fatal(srv.ListenAndServe()) 

	// Graceful Shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error!!", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit  // We're waiting

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown Failed: ", err)
	}
}
