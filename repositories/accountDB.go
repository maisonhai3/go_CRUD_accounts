package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
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

func (h *DBHandler) GetAll(ctx context.Context, limit int) ([]Account, error) {
	rows, err := h.DBConn.QueryContext(ctx, `
		SELECT id, name, currency, balance, created_at, updated_at
		FROM accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ?`, limit)
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

func (h *DBHandler) CreateAccount(ctx context.Context, acc CreateAccountParams) (string, error) {
	now := time.Now()
	var a = Account{
		Currency:  acc.Currency,
		Name:      acc.Name,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := h.DBConn.ExecContext(ctx, `INSERT INTO accounts (currency, name, balance, created_at) VALUES (?,?,?,?)`,
		a.Currency, a.Name, a.Balance, a.CreatedAt)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}
