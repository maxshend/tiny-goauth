package auth

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/maxshend/tiny_goauth/authtest"
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
		token := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodHS256, claims)

		if _, err := ValidateAccessToken(token); err != nil {
			t.Errorf("unexpected error: %q", err)
		}
	})

	t.Run("with invalid token", func(t *testing.T) {
		expired := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodHS256, expiredClaims)
		invalidSign := authtest.GenerateFakeJWT(t, []byte("invalid"), jwt.SigningMethodHS256, claims)
		invalidAlg := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodHS512, claims)

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
