package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	dbrepository "phase_1_foundations/db_repository"
	"time"
)

func (http_handler *HTTPHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Read max bytes (because this is a POST request)
	body := http.MaxBytesReader(w, r.Body, 1<<10) // 1 KiB

	// Unmarshal it
	var requestedAcc CreateAccountRequestDTO
	rDecoder := json.NewDecoder(body)
	rDecoder.DisallowUnknownFields() // Drop all fields that not in the contract
	if err := rDecoder.Decode(&requestedAcc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validations
	if !isValidCurrency(requestedAcc.Currency) {
		http.Error(w, "Invalid currency", http.StatusBadRequest)
		return
	}

	// ---- Boudary of Trust ----
	// Convert it to a DB DTO
	var newAcc = dbrepository.NewAccount{
		Name:     requestedAcc.Name,
		Currency: requestedAcc.Currency,
	}

	// Call DB to create
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	a, err := http_handler.DBRepo.CreateAccount(ctx, &newAcc)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	// Write to the response
	buf, _ := json.Marshal(a)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)

}
