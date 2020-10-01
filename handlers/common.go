package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/db"
	"github.com/maxshend/tiny_goauth/logwrapper"
)

// Deps contains dependencies of the http handlers
type Deps struct {
	DB         db.DataLayer
	Validator  *validator.Validate
	Translator ut.Translator
	Logger     *logwrapper.StandardLogger
}

type contextKey int

const contentTypeHeader = "Content-Type"
const jsonContentType = "application/json"
const successResponse = `{"success": true}`
const (
	tokenClaimsKey contextKey = iota
)

// Logout invalidates current JWT token
func Logout(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(deleteHandler(authenticatedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context().Value(tokenClaimsKey)
		claims, ok := c.(*auth.Claims)
		if !ok {
			log.Println("Cannot extract token claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		del, err := deps.DB.DeleteCache(claims.UUID)
		if err != nil || del == 0 {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		respondSuccess(w, http.StatusOK, nil)
	})))))
}

// Refresh returns a new access token if refresh token is valid
func Refresh(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		c, err := auth.ValidateRefreshToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := c.(*auth.Claims)
		if !ok {
			log.Println("Cannot extract token claims")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		del, err := deps.DB.DeleteCache(claims.UUID)
		if err != nil || del == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		td, err := auth.Token(claims.UserID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = saveTokenDetails(deps, claims.UserID, td)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		respondSuccess(w, http.StatusOK, td)
	}))))
}

func respondSuccess(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(status)
	if payload != nil {
		w.Write([]byte(response))
	}
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

func saveTokenDetails(deps *Deps, userID int64, td *auth.TokenDetails) error {
	at := time.Unix(td.AccessExpiresAt, 0)
	rt := time.Unix(td.RefreshExpiresAt, 0)
	now := time.Now()

	err := deps.DB.StoreCache(td.AccessUUID, userID, at.Sub(now))
	if err != nil {
		return err
	}

	err = deps.DB.StoreCache(td.RefreshUUID, userID, rt.Sub(now))
	if err != nil {
		return err
	}

	return nil
}
