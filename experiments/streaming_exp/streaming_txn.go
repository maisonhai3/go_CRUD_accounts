package streaming_exp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StreamingTxn(api string) {
	res, err := http.Get(api)
	if err != nil {
		log.Fatal() // enhance later
		return
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	// Opening the stream
	if _, err := decoder.Token(); err != nil {
		log.Fatal()
		return
	}

	for decoder.More() {
		var txn Transaction

		if err := decoder.Decode(&txn); err != nil {
			log.Fatal()
		}

		fmt.Printf("ID: %s", txn.ID)
		fmt.Printf("Amount: %f", txn.Amount)
		fmt.Printf("Status: %s", txn.Status)

	}

	// Closing the stream
	if _, err := decoder.Token(); err != nil {
		log.Fatal()
		return
	}
}
