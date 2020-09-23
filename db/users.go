package db

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/maxshend/tiny_goauth/models"
)

// CreateUser creates a new record in users database table
func (store *datastore) CreateUser(user *models.User) error {
	err := store.pool.QueryRow(
		context.Background(),
		"INSERT INTO users(email, password) VALUES($1, $2) RETURNING id, created_at",
		user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (store *datastore) ExistsWithField(fl validator.FieldLevel) (bool, error) {
	result := false
	err := store.pool.QueryRow(
		context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE "+fl.FieldName()+" = $1)", fl.Field().String(),
	).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}
