package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maxshend/tiny_goauth/auth"
)

// User represents data of a user in users table
type User struct {
	ID        int       `db:"id" json:"id"`
	Login     string    `db:"login" json:"login"`
	Password  string    `db:"password" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// EncryptPassword encrypts users password
func (user *User) EncryptPassword() error {
	hash, err := auth.EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	return nil
}

// CreateUser creates a new record in users database table
func CreateUser(db *pgxpool.Pool, user *User) error {
	row := db.QueryRow(
		context.Background(),
		"INSERT INTO users(login, password) VALUES($1, $2) RETURNING id, created_at",
		user.Login, user.Password,
	)

	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
