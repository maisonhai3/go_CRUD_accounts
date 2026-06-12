package handlers

import "time"

// CreateAccountRequest is the JSON body a client sends to create an account.
type CreateAccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

// ListAccountsQuery holds the validated query params for the list endpoint.
type ListAccountsQuery struct {
	Currency string
	Limit    int
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
