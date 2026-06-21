package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"accountCRUD/repositories"
)

func (h *Handler) GetAccountById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	if err := ctx.Err(); err != nil {
		return // The client gone, forget the work.
	}

	// Real DB call
	id := r.PathValue("id")
	a, err := h.Repo.GetById(ctx, id)

	if errors.Is(err, repositories.ErrAccNotFound) {
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

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// I DO NOT believe you, `r`
	rByte := http.MaxBytesReader(w, r.Body, 1<<20)
	rDecoded := json.NewDecoder(rByte)

	rDecoded.DisallowUnknownFields()

	var createAccRq CreateAccountRequest
	err := rDecoded.Decode(&createAccRq)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Looks good, but I must be in doubt
	if createAccRq.Name == "" || len(createAccRq.Name) > 300 {
		http.Error(w, "name is too long", http.StatusBadRequest)
		return
	}

	if !isValidCurrency(createAccRq.Currency) {
		http.Error(w, "currency is not supported", http.StatusBadRequest)
		return
	}

	// ----- Commit point
	var params = repositories.CreateAccountParams{
		Name:     createAccRq.Name,
		Currency: createAccRq.Currency,
	}
	newAcc, err := h.Repo.CreateAccount(ctx, params)
	if err != nil {
		http.Error(w, "unable to create new account", http.StatusInternalServerError)
		return
	}

	buf, _ := json.Marshal(toAccountResponse(newAcc))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Trusted Boundary mindset

	// --------------------
	// I don't believe you, requesters

	// Don't parse the whole request.
	// Trunk it
	if len(r.URL.RawQuery) > 100 {
		http.Error(w, "query is too long", http.StatusBadRequest)
		return
	}

	// Drop disallowed fields using what?
	query := r.URL.Query()
	allowed := map[string]bool{
		"currency": true,
		"limit":    true,
	}

	for key := range query {
		if !allowed[key] {
			http.Error(w, "key is not allowed", http.StatusBadRequest)
			return
		}
	}

	if len(query["currency"]) > 1 {
		http.Error(w, "many currencies is not allowed", http.StatusBadRequest)
		return
	}

	// Okay, you provide rightful fields, but not enough
	// Validations
	currency := query.Get("currency")
	if currency != "" && !isValidCurrency(currency) {
		http.Error(w, "currency is invalid", http.StatusBadRequest)
		return
	}

	limit := 10
	limitStr := query.Get("limit")
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil || parsed <= 0 || parsed > 100 {
			http.Error(w, "limit is invalid", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	// OK, I trust you.
	// Now, I will invoke DB to get your data

	// -------------------- 2
	// DBMS
	//a, err := h.DBConn.ExecContext(ctx, `GET * FROM accounts as a WHERE a.currency IS ? LIMIT ?`, currency, limit)
	accounts, err := h.Repo.GetAll(ctx, limit, currency)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// DB, you give me some data, but I do not trust you.
	// So, I encode it into DTO by myself
	getAccounts := toAccountResponses(accounts)

	// Okay, We've done with data extraction

	// -------------------- 3
	// Write it back
	buf, err := json.Marshal(getAccounts)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
}

func toAccountResponse(a repositories.Account) AccountResponse {
	return AccountResponse{
		ID:        a.ID,
		Name:      a.Name,
		Currency:  a.Currency,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func toAccountResponses(aS []repositories.Account) []AccountResponse {
	var out = []AccountResponse{}

	for _, a := range aS {
		newA := toAccountResponse(a)
		out = append(out, newA)
	}

	return out
}
