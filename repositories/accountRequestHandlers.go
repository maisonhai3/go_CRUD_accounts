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
