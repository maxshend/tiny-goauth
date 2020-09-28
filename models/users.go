package models

import (
	"time"
)

// User represents data of a user in users table
type User struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email" validate:"required,email,unique_user"`
	Password  string    `db:"password" json:"password" validate:"required,password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
