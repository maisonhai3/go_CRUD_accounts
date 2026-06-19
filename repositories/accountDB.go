package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
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

// rowScanner is satisfied by both *sql.Row (QueryRow) and *sql.Rows (Query),
// so scanAccount can serve the single-row and multi-row reads alike.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanAccount reads one account row. Timestamps are stored as RFC3339 TEXT,
// so they're scanned into strings and parsed back into time.Time.
func scanAccount(s rowScanner) (Account, error) {
	var a Account
	var createdAt, updatedAt string
	if err := s.Scan(&a.ID, &a.Name, &a.Currency, &a.Balance, &createdAt, &updatedAt); err != nil {
		return Account{}, err
	}

	var err error
	if a.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return Account{}, err
	}
	if a.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		return Account{}, err
	}

	return a, nil
}

func (h *DBHandler) GetById(ctx context.Context, id string) (Account, error) {
	// Emit a query to db with Scan()
	a, err := scanAccount(h.DBConn.QueryRowContext(ctx,
		`
			SELECT id, name, currency, balance, created_at, updated_at
			FROM accounts
			WHERE id = ? AND deleted_at IS NULL
		`, id))

	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, nil
	}

	return a, err
}

func (h *DBHandler) GetAll(ctx context.Context, limit int, currency string) ([]Account, error) {
	cond := []string{"deleted_at IS NULL"}
	args := []any{}

	if currency != "" {
		cond = append(cond, "currency = ?")
		args = append(args, currency)
	}

	query := `
		SELECT id, name, currency, balance, created_at, updated_at
		FROM accounts
		WHERE ` + strings.Join(cond, " AND ") + `
		ORDER BY created_at DESC
		LIMIT ?`
	args = append(args, limit)

	rows, err := h.DBConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Account
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}

		out = append(out, a)
	}

	return out, rows.Err()
}

func (h *DBHandler) CreateAccount(ctx context.Context, acc CreateAccountParams) (Account, error) {
	// Store timestamps as RFC3339 TEXT so they round-trip cleanly on read.
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	result, err := h.DBConn.ExecContext(ctx,
		`INSERT INTO accounts
			(currency, name, balance, created_at, updated_at)
			VALUES (?,?,?,?,?)`,
		acc.Currency, acc.Name, 0, nowStr, nowStr)
	if err != nil {
		return Account{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:        strconv.FormatInt(id, 10),
		Name:      acc.Name,
		Currency:  acc.Currency,
		Balance:   0,
		CreatedAt: now.Truncate(time.Second),
		UpdatedAt: now.Truncate(time.Second),
	}, nil
}
