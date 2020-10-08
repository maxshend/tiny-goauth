package auth

import (
	"testing"

	"github.com/maxshend/tiny_goauth/authtest"
)

func TestEncryptPassword(t *testing.T) {
	t.Run("creates nonempty hash string", func(t *testing.T) {
		result, err := EncryptPassword("12345678")

		if err != nil {
			t.Errorf("got an error %q", err)
		}

		if len(result) == 0 {
			t.Errorf("got an empty hash string")
		}
	})

	t.Run("returns error for an empty password", func(t *testing.T) {
		_, err := EncryptPassword("")

		if err == nil {
			t.Fatalf("should got an error")
		}

		authtest.AssertError(t, errEmptyPassword, err)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("validates that the hash generated from the password", func(t *testing.T) {
		pswd := "foobar"
		hash, _ := EncryptPassword(pswd)
		valid := ValidatePassword(pswd, hash)

		if !valid {
			t.Errorf("password should be valid")
		}
	})

	t.Run("returns an error for invalid hash/password combination", func(t *testing.T) {
		pswd := "foobar"
		hash, _ := EncryptPassword("invalid")
		valid := ValidatePassword(pswd, hash)

		if valid {
			t.Errorf("password should be invalid")
		}
	})
}
