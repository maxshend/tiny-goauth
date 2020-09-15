package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// User represents data of a user in users table
type User struct {
	ID        int       `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// CreateUser creates a new record in users database table
func CreateUser(db *pgxpool.Pool, user *User) error {
	// TODO: Hash password
	encryptedPassword := user.Password
	row := db.QueryRow(
		context.Background(),
		"INSERT INTO users(login, password) VALUES($1, $2) RETURNING id, created_at",
		user.Login, encryptedPassword,
	)

	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
