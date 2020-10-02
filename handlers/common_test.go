package handlers

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestLogout(t *testing.T) {
	t.Run("returns MethodNotAllowed for non-DELETE requests", func(t *testing.T) {
		recorder := performRequest(t, "GET", "/logout", Logout, nil, jsonHeaders)

		validateStatusCode(t, recorder, http.StatusMethodNotAllowed)
	})

	t.Run("returns BadRequest without json 'Conten-Type' header", func(t *testing.T) {
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, nil)

		validateStatusCode(t, recorder, http.StatusBadRequest)
	})

	secret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	token := generateFakeJWT(t, secret, jwt.SigningMethodHS256, claims)
	expired := generateFakeJWT(t, secret, jwt.SigningMethodHS256, expiredClaims)

	t.Run("returns Unauthorized with invalid 'Authorization' header", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = expired
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h)

		validateStatusCode(t, recorder, http.StatusUnauthorized)
	})

	t.Run("returns OK with valid token", func(t *testing.T) {
		h := jsonHeaders
		h[auhtorizationHeader] = token
		recorder := performRequest(t, "DELETE", "/logout", Logout, nil, h)

		validateStatusCode(t, recorder, http.StatusOK)
	})
}

func generateFakeJWT(t *testing.T, sign []byte, method jwt.SigningMethod, claims jwt.Claims) string {
	t.Helper()

	jwtToken := jwt.NewWithClaims(method, claims)
	token, err := jwtToken.SignedString(sign)
	if err != nil {
		t.Fatal(err)
	}

	return token
}
