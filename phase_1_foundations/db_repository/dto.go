package dbrepository

type NewAccount struct {
	Name     string
	Currency string
}

type CreatedAccount struct {
	ID   int64
	Name string
}