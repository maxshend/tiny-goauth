package auth

import "golang.org/x/crypto/bcrypt"

// EncryptPassword generates hash from a password string
func EncryptPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errEmptyPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

// ValidatePassword validates equality of a password hash and a password string
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
