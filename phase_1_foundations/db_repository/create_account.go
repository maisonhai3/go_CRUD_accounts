package dbrepository

import "context"

func (db *DBRepo) CreateAccount(ctx context.Context, a *NewAccount) (*CreatedAccount, error) {
	const balance = 0 // server always forces a new account to start at zero

	res, err := db.Conn.ExecContext(ctx,
		`INSERT INTO accounts (name, currency, balance) VALUES (?, ?, ?)`,
		a.Name, a.Currency, balance)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &CreatedAccount{ID: id, Name: a.Name}, nil
}
