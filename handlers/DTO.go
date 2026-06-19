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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
