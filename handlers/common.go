package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/maxshend/tiny_goauth/auth"
	"github.com/maxshend/tiny_goauth/db"
	"github.com/maxshend/tiny_goauth/logwrapper"
	"github.com/maxshend/tiny_goauth/models"
)

// Deps contains dependencies of the http handlers
type Deps struct {
	DB         db.DataLayer
	Validator  *validator.Validate
	Translator ut.Translator
	Logger     *logwrapper.StandardLogger
	Keys       *auth.RSAKeys
}

type contextKey int

const contentTypeHeader = "Content-Type"
const auhtorizationHeader = "Authorization"
const jsonContentType = "application/json"
const successResponse = `{"success": true}`
const invalidTokenMsg = "Invalid Authorization token."
const postMethod = "POST"
const (
	tokenClaimsKey contextKey = iota
)
const maxBodySize = 1048576
const defaultUsersEndpoint = "/internal/tiny_goauth/users"

type handlerErr string

func (e handlerErr) Error() string { return string(e) }

const failExternalResponse = handlerErr("External service returned invalid response")

// Logout invalidates current JWT token
func Logout(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(deleteHandler(authenticatedHandler(deps, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context().Value(tokenClaimsKey)
		claims, ok := c.(*auth.Claims)
		if !ok {
			respondInvalidToken(w)
			return
		}

		del, err := deps.DB.DeleteCache(claims.UUID)
		if del == 0 {
			respondInvalidToken(w)
			return
		}
		if err != nil {
			deps.Logger.RequestError(r, err)
			respondInternalError(w)
			return
		}

		respondSuccess(w, http.StatusOK, nil)
	})))))
}

// Refresh returns a new access token if refresh token is valid
func Refresh(deps *Deps) http.Handler {
	return logHandler(deps, jsonHandler(postHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		c, err := auth.ValidateToken(token, deps.Keys.RefreshVerify)
		if err != nil {
			respondInvalidToken(w)
			return
		}

		claims, ok := c.(*auth.Claims)
		if !ok {
			respondInvalidToken(w)
			return
		}

		del, err := deps.DB.DeleteCache(claims.UUID)
		if del == 0 {
			respondInvalidToken(w)
			return
		}
		if err != nil {
			deps.Logger.RequestError(r, err)
			respondInternalError(w)
			return
		}

		td, err := auth.Token(claims.UserID, claims.Roles, deps.Keys)
		if err != nil {
			respondInvalidToken(w)
			return
		}

		err = saveTokenDetails(deps, claims.UserID, td)
		if err != nil {
			respondInvalidToken(w)
			return
		}

		payload, err := json.Marshal(td)
		if err != nil {
			deps.Logger.RequestError(r, err)
			respondInternalError(w)
			return
		}

		respondSuccess(w, http.StatusOK, payload)
	}))))
}

func respondSuccess(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondInternalError(w)
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

func respondInvalidToken(w http.ResponseWriter) {
	respondError(w, http.StatusUnauthorized, invalidTokenMsg)
}

func respondInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
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

func postJSON(url string, body io.Reader) (result io.Reader, code int, err error) {
	code = http.StatusServiceUnavailable

	req, err := http.NewRequest(postMethod, url, body)
	if err != nil {
		return nil, code, err
	}

	req.Header.Set(contentTypeHeader, jsonContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, code, err
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	respBody, _ := ioutil.ReadAll(resp.Body)

	return bytes.NewReader(respBody), code, nil
}

func createExternalUser(user *models.User) error {
	endpoint, found := os.LookupEnv("API_USERS_ENDPOINT")
	if !found || len(endpoint) == 0 {
		endpoint = defaultUsersEndpoint
	}

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, code, err := postJSON(os.Getenv("API_HOST")+endpoint, bytes.NewReader(body))

	if code != http.StatusOK {
		return failExternalResponse
	}

	return err
}
