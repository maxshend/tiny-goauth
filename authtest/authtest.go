package authtest

import (
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// GenerateFakeJWT generates fake JWT token with valid format
func GenerateFakeJWT(t *testing.T, sign []byte, method jwt.SigningMethod, claims jwt.Claims) string {
	t.Helper()

	jwtToken := jwt.NewWithClaims(method, claims)
	token, err := jwtToken.SignedString(sign)
	if err != nil {
		t.Fatal(err)
	}

	return token
}

// AssertStatusCode asserts HTTP status code
func AssertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if got := recorder.Code; got != expected {
		t.Errorf("Returned wrong status code. Expected %d, got %d", expected, got)
	}
}
