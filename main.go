package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()
}
