package dbrepository

import "database/sql"

type DBRepo struct {
	Conn *sql.DB
}
