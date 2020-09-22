package models

import (
	"time"

	"github.com/maxshend/tiny_goauth/auth"
)

// User represents data of a user in users table
type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email" validate:"required,email,unique_user"`
	Password  string    `db:"password" json:"password" validate:"required,password"`
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
