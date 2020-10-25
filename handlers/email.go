package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/models"
)

// EmailRegister handles email registration requests
func EmailRegister(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&user)
		if err != nil {
			deps.Logger.RequestError(r, err)
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
			deps.Logger.RequestError(r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Password = hash

		err = deps.DB.CreateUser(&user)
		if err != nil {
			deps.Logger.RequestError(r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = createExternalUser(&user)
		if err != nil {
			deps.Logger.RequestError(r, err)

			err = deps.DB.DeleteUser(user.ID)
			if err != nil {
				deps.Logger.RequestError(r, err)
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := auth.Token(user.ID, user.Roles, deps.Keys)
		if err != nil {
			deps.Logger.RequestError(r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = saveTokenDetails(deps, user.ID, token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondSuccess(w, http.StatusOK, token)
	}))))
}

// EmailLogin validates user email and password combination
func EmailLogin(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var loginUser models.User
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&loginUser)
		if err != nil {
			deps.Logger.RequestError(r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := deps.DB.UserByEmail(loginUser.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !auth.ValidatePassword(loginUser.Password, user.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := auth.Token(user.ID, user.Roles, deps.Keys)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		err = saveTokenDetails(deps, user.ID, token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondSuccess(w, http.StatusOK, token)
	}))))
}
