package authtest

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

// GenerateFakeJWT generates fake JWT token with valid format
func GenerateFakeJWT(t *testing.T, sign interface{}, method jwt.SigningMethod, claims jwt.Claims) string {
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

// AssertError asserts error type
func AssertError(t *testing.T, want, got error) {
	t.Helper()

	if want != got {
		t.Errorf("expected %q got %q", want, got)
	}
}

// GeneratePrivateKey creates a RSA Private Key
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 768)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// GeneratePublicKey creates a RSA Public Key
func GeneratePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(publicKey), nil
}
