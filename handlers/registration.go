package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/models"
)

// EmailRegister handles email registration requests
func EmailRegister(deps *Deps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get(contentTypeHeader) != jsonContentType {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&user)
		if err != nil {
			http.Error(w, "Error while decoding json body", http.StatusInternalServerError)
			return
		}

		err = deps.Validator.Struct(&user)
		if err != nil {
			respondModelError(deps, w, err.(validator.ValidationErrors))
			return
		}

		err = user.EncryptPassword()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = deps.DB.CreateUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respondSuccess(w, http.StatusOK, &user)

		return
	}
}
