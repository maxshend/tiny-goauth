package models

import (
	"testing"

	"github.com/maxshend/tiny_goauth/validations"
)

type TestData struct {
	User  *User
	Valid bool
}

func TestUser(t *testing.T) {
	testData := []TestData{
		TestData{User: &User{Email: "invalid.mail.com", Password: "12345678"}, Valid: false},
		TestData{User: &User{Email: "valid@mail.com", Password: "12345678"}, Valid: true},
		TestData{User: &User{Email: "valid@mail.com", Password: "1234567"}, Valid: false},
	}

	validate, _, _ := validations.Init()

	for _, data := range testData {
		err := validate.Struct(data.User)
		if err != nil && data.Valid {
			t.Errorf("Exptected  %+v to be valid but was invalid instead.", data.User)
		} else if err == nil && !data.Valid {
			t.Errorf("Exptected  %+v to be invalid but was valid instead.", data.User)
		}
	}
}
