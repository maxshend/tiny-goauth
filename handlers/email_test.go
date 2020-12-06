package handlers

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/maxshend/tiny_goauth/authtest"
)

var jsonHeaders = map[string]string{contentTypeHeader: jsonContentType}

func TestEmailRegister(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/register", EmailRegister, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns UnprocessableEntity when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, strings.NewReader("invalid"), jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns UnprocessableEntity with invalid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnprocessableEntity)
	})

	t.Run("returns OK with valid user data", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "valid@mail.com", "password": "12345678"}`))

		externalApp := testServer()
		defer externalApp.Close()

		os.Setenv("API_HOST", externalApp.URL)

		recorder := performRequest(t, "POST", "/email/register", EmailRegister, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func TestEmailLogin(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/email/login", EmailLogin, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Content-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	t.Run("returns InternalServerError when body isn't valid json", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, strings.NewReader("invalid"), jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusInternalServerError)
	})

	t.Run("returns Unauthorized with invalid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "invalid.mail.com", "password": "foobar123"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid user creds", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{"email": "test@mail.com", "password": "password"}`))
		recorder := performRequest(t, "POST", "/email/login", EmailLogin, body, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}
