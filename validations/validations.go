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
func Init() (*validator.Validate, ut.Translator, error) {
	en := en.New()
	uni = ut.New(en, en)

	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, nil, errors.New("Translator cannot be located")
	}

	v := validator.New()

	en_translations.RegisterDefaultTranslations(v, trans)

	err := v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	if err != nil {
		return nil, nil, err
	}

	err = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} has invalid format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})
	if err != nil {
		return nil, nil, err
	}

	err = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("passwd", "{0} minimal length is 8 characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})
	if err != nil {
		return nil, nil, err
	}

	err = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 7
	})
	if err != nil {
		return nil, nil, err
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v, trans, nil
}
