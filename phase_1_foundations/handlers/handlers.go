package handlers

import "accountCRUD/repositories"

// Handler holds the dependencies for the HTTP transport layer.
// It depends on the data layer rather than being it.
type Handler struct {
	Repo *repositories.DBHandler
}
