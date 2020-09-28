package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/db"
)

// Deps contains dependencies of the http handlers
type Deps struct {
	DB         db.DataLayer
	Validator  *validator.Validate
	Translator ut.Translator
}

type contextKey int

const contentTypeHeader = "Content-Type"
const jsonContentType = "application/json"
const (
	tokenClaimsKey contextKey = iota
)
const successResponse = `{"success": true}`

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

func authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		claims, err := auth.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), tokenClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jsonOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(contentTypeHeader) != jsonContentType {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
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

// Logout invalidates current JWT token
func Logout(deps *Deps) http.Handler {
	return jsonOnly(authenticatedOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

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
	})))
}
