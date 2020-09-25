package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/auth"
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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = deps.Validator.Struct(&user)
		if err != nil {
			respondModelError(deps, w, err.(validator.ValidationErrors))
			return
		}

		hash, err := auth.EncryptPassword(user.Password)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Password = hash

		err = deps.DB.CreateUser(&user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := auth.Token(user.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondSuccess(w, http.StatusOK, token)

		return
	}
}

// EmailLogin validates user email and password combination
func EmailLogin(deps *Deps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get(contentTypeHeader) != jsonContentType {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var loginUser models.User
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&loginUser)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := deps.DB.UserByEmail(loginUser.Email)
		if err != nil {
			log.Printf("Error while getting user by email: %q\n", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := auth.Token(user.ID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !auth.ValidatePassword(loginUser.Password, user.Password) {
			log.Printf("Invalid login credentials: %+v\n", loginUser)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		respondSuccess(w, http.StatusOK, token)

		return
	}
}
