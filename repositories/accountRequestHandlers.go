package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (h *DBHandler) GetAccountById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	if err := ctx.Err(); err != nil {
		return // The client gone, forget the work.
	}

	// Real DB call
	id := r.PathValue("id")
	a, err := h.GetById(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// DB timeout
	if err := ctx.Err(); err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	// Serialize
	j, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Error in marshaling", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(j)
}

func (h *DBHandler) GetAccountByUsername(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// I DO NOT believe you, `r`
	rByte := http.MaxBytesReader(w, r.Body, 1<<20)
	rDecoded := json.NewDecoder(rByte)

	rDecoded.DisallowUnknownFields()

	var createAccRq createAccSchema
	err := rDecoded.Decode(&createAccRq)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Looks good, but I must be in doubt
	if createAccRq.name == "" || len(createAccRq.name) > 30 {
		http.Error(w, "name is too long", http.StatusBadRequest)
		return
	}

	if !isValidCurrency(createAccRq.currency) {
		http.Error(w, "currency is not supported", http.StatusBadRequest)
		return
	}

	// ----- Commit point
	var accDTO = AccountDTO{
		ID:       "123abc",
		Name:     createAccRq.name,
		Currency: createAccRq.currency,
	}
	id, err := h.createAccount(ctx, accDTO)
	if err != nil {
		http.Error(w, "unable to create new account", http.StatusInternalServerError)
		return
	}

	buf, _ := json.Marshal(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

type createAccSchema struct {
	name     string
	currency string
}

func isValidCurrency(cur string) bool {
	// Whitelist
	switch cur {
	case "VND", "USD":
		return true
	}
	return false
}

type AccountDTO struct {
	ID       string
	Name     string
	Currency string
}
