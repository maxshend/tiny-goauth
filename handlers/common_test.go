package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/maxshend/tiny_goauth/authtest"
)

func TestLogout(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-DELETE requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/logout", Logout, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	privateKey, err := authtest.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	token := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)
	expired := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, expiredClaims)

	t.Run("returns Unauthorized with invalid 'Authorization' header", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = expired
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid token", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = token
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}

func TestRefresh(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-POST requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/refresh", Refresh, nil, jsonHeaders, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, nil, nil)

		authtest.AssertStatusCode(t, recorder, http.StatusBadRequest)
	})

	privateKey, err := authtest.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	token := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)
	expired := authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, expiredClaims)

	t.Run("returns Unauthorized with invalid 'Authorization' header", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = expired
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusUnauthorized)
	})

	token = authtest.GenerateFakeJWT(t, privateKey, jwt.SigningMethodRS256, claims)

	t.Run("returns OK with valid Refresh token", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = token
		recorder := performRequest(t, "POST", "/refresh", Refresh, nil, h, privateKey)

		authtest.AssertStatusCode(t, recorder, http.StatusOK)
	})
}
