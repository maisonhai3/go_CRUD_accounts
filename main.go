package main

import (
	"fmt"
	"log"
	"net/http"
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

	fmt.Fprintf(w, "Write")
}

func main() {
	http.HandleFunc("/accounts", getAccounts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
