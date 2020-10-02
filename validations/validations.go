package validations

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/maxshend/tiny_goauth/db"
)

var uni *ut.UniversalTranslator
var errCannotLocateTranslator = errors.New("Translator cannot be located")

// Init creates new validator instance
func Init(db db.DataLayer) (validate *validator.Validate, translator ut.Translator, err error) {
	en := en.New()
	uni = ut.New(en, en)

	translator, found := uni.GetTranslator("en")
	if !found {
		err = errCannotLocateTranslator
		return
	}

	validate = validator.New()

	en_translations.RegisterDefaultTranslations(validate, translator)

	err = registerTranslation(validate, translator, "required", "{0} is a required field")
	if err != nil {
		return
	}

	err = registerTranslation(validate, translator, "email", "{0} has invalid format")
	if err != nil {
		return
	}

	err = registerTranslation(validate, translator, "password", "{0} minimal length is 8 characters")
	if err != nil {
		return
	}

	err = registerTranslation(validate, translator, "unique_user", "this {0} is already taken")
	if err != nil {
		return
	}

	err = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 7
	})
	if err != nil {
		return
	}

	err = validate.RegisterValidation("unique_user", func(fl validator.FieldLevel) bool {
		exists, err := db.UserExistsWithField(fl)
		if err != nil {
			return false
		}

		return !exists
	})
	if err != nil {
		return
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return
}

func registerTranslation(validate *validator.Validate, translator ut.Translator, name, message string) error {
	return validate.RegisterTranslation(name, translator, func(ut ut.Translator) error {
		return ut.Add(name, message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(name, fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
}
