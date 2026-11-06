package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(dsn string) (*sql.DB, error) {
	db, err	:= sql.Open("postgres", dsn)
	if err != nil {
		return  nil , fmt.Errorf("failed to open db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	/*
		pingErr := db.Ping()
		if pingErr != nil {
			return nil, fmt.Errorf("failed to ping db: %w", err)
		}
	*/

	return db, nil;
}