package repositories

import "time"

type CreateAccountDTO struct {
	ID       string
	Name     string
	Currency string
}

type GetAccountDTO struct {
	ID        string
	Name      string
	Currency  string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
