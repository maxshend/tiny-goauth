package handlers

import (
	"encoding/json"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/db"
)

// Deps contains dependencies of the http handlers
type Deps struct {
	DB         db.DataLayer
	Validator  *validator.Validate
	Translator ut.Translator
}

const contentTypeHeader = "Content-Type"
const jsonContentType = "application/json"

func respondSuccess(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message interface{}) {
	respondSuccess(w, code, map[string]interface{}{"errors": message})
}

func respondModelError(deps *Deps, w http.ResponseWriter, err validator.ValidationErrors) {
	errResponse := make(map[string]string)
	for _, err := range err {
		errResponse[err.Field()] = err.Translate(deps.Translator)
	}

	respondError(w, http.StatusUnprocessableEntity, errResponse)
}
