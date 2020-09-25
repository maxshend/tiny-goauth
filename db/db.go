package db

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxshend/tiny_goauth/models"
)

// DataLayer is the interface that wraps methods to access database
type DataLayer interface {
	CreateUser(*models.User) error
	UserExistsWithField(fl validator.FieldLevel) (bool, error)
	UserByEmail(string) (*models.User, error)
	Close()
}

type datastore struct {
	pool *pgxpool.Pool
}

// Init initializes connection to the database
func Init(dataURL string) (DataLayer, error) {
	pool, err := pgxpool.Connect(context.Background(), dataURL)
	if err != nil {
		return nil, err
	}

	return &datastore{pool: pool}, nil
}

// Close closes connection to the database
func (store *datastore) Close() {
	store.pool.Close()
}
