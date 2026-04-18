package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns, maxIddleConns int, maxIddleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIddleConns)

	duration, err := time.ParseDuration(maxIddleTime)
	if err != nil {
		return nil, err
	}
	
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()


	if err :=db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}