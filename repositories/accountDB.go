package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Account struct {
	ID        string
	Name      string
	Currency  string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// ----------- DB Handlers ------------

func (h *DBHandler) GetById(ctx context.Context, id string) (Account, error) {
	var a Account
	// Emit a query to db with Scan()
	err := h.DBConn.QueryRowContext(ctx,
		`
			SELECT id, name, currency, balance, created_at, updated_at
			FROM accounts
			WHERE id = ? AND deleted_at IS NULL
		`, id).Scan(&a.ID, &a.Name, &a.Currency, &a.Balance, &a.CreatedAt, &a.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, err
	}

	return a, err
}

func (h *DBHandler) getAll(ctx context.Context) ([]Account, error) {
	rows, err := h.DBConn.QueryContext(ctx, `
		SELECT id, name, currency, balance, created_at, updated_at
		FROM accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Currency, &a.Balance, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}

		out = append(out, a)
	}

	return out, rows.Err()
}

// ----------- Request Handlers ------------
