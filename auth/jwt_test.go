package auth

import (
	"log"
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

func generateFakeJWT(exp int64, sign []byte, method jwt.SigningMethod) string {
	jwtToken := jwt.NewWithClaims(method, jwt.MapClaims{
		"exp": exp,
	})
	token, err := jwtToken.SignedString(sign)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func TestValidateToken(t *testing.T) {
	secret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	validTime := time.Now().Add(time.Minute * 15).Unix()

	t.Run("with valid token", func(t *testing.T) {
		token := generateFakeJWT(validTime, secret, jwt.SigningMethodHS256)

		if _, err := ValidateToken(token); err != nil {
			t.Errorf("unexpected error: %q", err)
		}
	})

	t.Run("with invalid token", func(t *testing.T) {
		expired := generateFakeJWT(time.Now().Add(-(time.Minute * 15)).Unix(), secret, jwt.SigningMethodHS256)
		invalidSign := generateFakeJWT(validTime, []byte("invalid"), jwt.SigningMethodHS256)
		invalidAlg := generateFakeJWT(validTime, secret, jwt.SigningMethodHS512)

		tokenCases := []struct {
			title string
			token string
			msg   string
		}{
			{title: "Expired", token: expired, msg: "Token is expired"},
			{title: "Invalid signature", token: invalidSign, msg: "signature is invalid"},
			{title: "Invalid format", token: "foobar", msg: "token contains an invalid number of segments"},
			{title: "Invalid signing method", token: invalidAlg, msg: "Unexpected signing method: HS512"},
		}

		for _, tc := range tokenCases {
			t.Run(tc.title, func(t *testing.T) {
				_, err := ValidateToken(tc.token)

				if err == nil {
					t.Errorf("expected to be invalid")
				} else if err.Error() != tc.msg {
					t.Errorf("expected %q got %q", tc.msg, err)
				}
			})
		}
	})
}
