package auth

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestToken(t *testing.T) {
	t.Run("without errors", func(t *testing.T) {
		_, err := Token(0)
		if err != nil {
			t.Errorf("got %q error", err.Error())
		}
	})

	t.Run("returns non empty tokens", func(t *testing.T) {
		details, _ := Token(0)
		if len(details.Access) == 0 || len(details.Refresh) == 0 {
			t.Error("got empty tokens")
		}
	})
}

func TestValidateAccessToken(t *testing.T) {
	secret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}

	t.Run("with valid token", func(t *testing.T) {
		token := generateFakeJWT(t, secret, jwt.SigningMethodHS256, claims)

		if _, err := ValidateAccessToken(token); err != nil {
			t.Errorf("unexpected error: %q", err)
		}
	})

	t.Run("with invalid token", func(t *testing.T) {
		expired := generateFakeJWT(t, secret, jwt.SigningMethodHS256, expiredClaims)
		invalidSign := generateFakeJWT(t, []byte("invalid"), jwt.SigningMethodHS256, claims)
		invalidAlg := generateFakeJWT(t, secret, jwt.SigningMethodHS512, claims)

		tokenCases := []struct {
			title string
			token string
			msg   string
		}{
			{title: "Expired", token: expired, msg: "token is expired by 15m0s"},
			{title: "Invalid signature", token: invalidSign, msg: "signature is invalid"},
			{title: "Invalid format", token: "foobar", msg: "token contains an invalid number of segments"},
			{title: "Invalid signing method", token: invalidAlg, msg: "Unexpected signing method: HS512"},
		}

		for _, tc := range tokenCases {
			t.Run(tc.title, func(t *testing.T) {
				_, err := ValidateAccessToken(tc.token)

				if err == nil {
					t.Fatal("expected to be invalid")
				}

				if err.Error() != tc.msg {
					t.Errorf("expected %q got %q", tc.msg, err)
				}
			})
		}
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
