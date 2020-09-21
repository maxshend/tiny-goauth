package validations

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var uni *ut.UniversalTranslator

// Init creates new validator instance
func Init() (validate *validator.Validate, translator ut.Translator, err error) {
	en := en.New()
	uni = ut.New(en, en)

	translator, found := uni.GetTranslator("en")
	if !found {
		err = errors.New("Translator cannot be located")
		return
	}

	validate = validator.New()

	en_translations.RegisterDefaultTranslations(validate, translator)

	err = validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("required", fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
	if err != nil {
		return
	}

	err = validate.RegisterTranslation("email", translator, func(ut ut.Translator) error {
		return ut.Add("email", "{0} has invalid format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("email", fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
	if err != nil {
		return
	}

	err = validate.RegisterTranslation("password", translator, func(ut ut.Translator) error {
		return ut.Add("password", "{0} minimal length is 8 characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T("password", fe.Field())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})
	if err != nil {
		return
	}

	err = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 7
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
