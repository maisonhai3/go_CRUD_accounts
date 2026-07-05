package repositories

// CreateAccountParams is the validated input passed to the DB layer when
// creating an account.
type CreateAccountParams struct {
	Name     string
	Currency string
}
