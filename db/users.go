package db

import (
	"github.com/go-playground/validator"
	"github.com/maxshend/tiny_goauth/models"
)

// CreateUser creates a new record in users database table
func (store *datastore) CreateUser(user *models.User) error {
	err := store.db.QueryRow(
		ctx,
		"INSERT INTO users(email, password) VALUES($1, $2) RETURNING id, created_at",
		user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (store *datastore) UserExistsWithField(fl validator.FieldLevel) (bool, error) {
	result := false
	err := store.db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE "+fl.FieldName()+" = $1)", fl.Field().String(),
	).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (store *datastore) UserByEmail(email string) (*models.User, error) {
	var user models.User
	err := store.db.QueryRow(
		ctx,
		"SELECT id, email, password, created_at FROM users WHERE email = $1 LIMIT 1", email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
