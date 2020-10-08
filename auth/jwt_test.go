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
		_, err := Token(0, nil)
		if err != nil {
			t.Errorf("got %q error", err.Error())
		}
	})

	t.Run("returns non empty tokens", func(t *testing.T) {
		details, _ := Token(0, nil)
		if details == nil || len(details.Access) == 0 || len(details.Refresh) == 0 {
			t.Error("got empty tokens")
		}
	})

	t.Run("returns error when refresh token isn't set", func(t *testing.T) {
		os.Unsetenv("ACCESS_TOKEN_SECRET")

		_, err := Token(0, nil)
		if err == nil {
			t.Errorf("should got an error")
		}
	})
}

func TestValidateAccessToken(t *testing.T) {
	testToken(t, []byte(os.Getenv("ACCESS_TOKEN_SECRET")), ValidateAccessToken)
}

func TestValidateRefreshToken(t *testing.T) {
	testToken(t, []byte(os.Getenv("REFRESH_TOKEN_SECRET")), ValidateRefreshToken)
}

func testToken(t *testing.T, secret []byte, fn func(tokenString string) (jwt.Claims, error)) {
	t.Helper()

	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}

	t.Run("with valid token", func(t *testing.T) {
		token := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodHS256, claims)

		if _, err := fn(token); err != nil {
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
				_, err := fn(tc.token)

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
