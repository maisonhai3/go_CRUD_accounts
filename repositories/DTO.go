package repositories

import "time"

// CreateAccountParams is the validated input passed to the DB layer when
// creating an account.
type CreateAccountParams struct {
	ID       string
	Name     string
	Currency string
}

// AccountResponse is the JSON shape returned to clients for an account.
type AccountResponse struct {
	ID        string
	Name      string
	Currency  string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
