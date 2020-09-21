package db

import (
	"context"

	"github.com/maxshend/tiny_goauth/models"
)

// CreateUser creates a new record in users database table
func (store *datastore) CreateUser(user *models.User) error {
	row := store.pool.QueryRow(
		context.Background(),
		"INSERT INTO users(email, password) VALUES($1, $2) RETURNING id, created_at",
		user.Email, user.Password,
	)

	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
