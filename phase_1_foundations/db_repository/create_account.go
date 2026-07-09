package dbrepository

import (
	"context"
	"time"
)

func (db *DBRepo) CreateAccount(ctx context.Context, a *NewAccount)(any, error){
	var nowForResponse = time.Now().UTC()
	var nowForDB = nowForResponse.Format(time.RFC3339)

	new, err := db.Conn.ExecContext(ctx, 
		`INSERT INTO accounts ()`,
		a.Name, nowForDB)
	if err  != nil	{
		return nil, err
	}
	
	return new, nil
}