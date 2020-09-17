package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitDB initializes connection to the database
func Init(dataURL string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), dataURL)
}
