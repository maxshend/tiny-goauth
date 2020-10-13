package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/maxshend/tiny_goauth/authtest"
)

func TestToken(t *testing.T) {
	privateKey, err := authtest.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	keys := &RSAKeys{AccessSign: privateKey, RefreshSign: privateKey}

	t.Run("without errors", func(t *testing.T) {
		_, err := Token(0, nil, keys)
		if err != nil {
			t.Errorf("got %q error", err.Error())
		}
	})

	t.Run("returns non empty tokens", func(t *testing.T) {
		details, _ := Token(0, nil, keys)
		if details == nil || len(details.Access) == 0 || len(details.Refresh) == 0 {
			t.Error("got empty tokens")
		}
	})
}

func TestValidateToken(t *testing.T) {
	accessSign, _ := authtest.GeneratePrivateKey()
	refreshSign, _ := authtest.GeneratePrivateKey()
	keys := &RSAKeys{AccessSign: accessSign, AccessVerify: &accessSign.PublicKey, RefreshSign: refreshSign}

	claims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * 15).Unix()}
	expiredClaims := jwt.MapClaims{"exp": time.Now().Add(time.Minute * -15).Unix()}
	secret := keys.AccessSign

	t.Run("with valid token", func(t *testing.T) {
		token := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodRS256, claims)

		if _, err := ValidateToken(token, keys.AccessVerify); err != nil {
			t.Errorf("unexpected error: %q", err)
		}
	})

	t.Run("with invalid token", func(t *testing.T) {
		expired := authtest.GenerateFakeJWT(t, secret, jwt.SigningMethodRS256, expiredClaims)
		invalidSign := authtest.GenerateFakeJWT(t, keys.RefreshSign, jwt.SigningMethodRS256, claims)
		invalidAlg := authtest.GenerateFakeJWT(t, []byte("foobar123"), jwt.SigningMethodHS512, claims)

		tokenCases := []struct {
			title string
			token string
			msg   string
		}{
			{title: "Expired", token: expired, msg: "token is expired by 15m0s"},
			{title: "Invalid signature", token: invalidSign, msg: "crypto/rsa: verification error"},
			{title: "Invalid format", token: "foobar", msg: "token contains an invalid number of segments"},
			{title: "Invalid signing method", token: invalidAlg, msg: "Unexpected signing method: HS512"},
		}

		for _, tc := range tokenCases {
			t.Run(tc.title, func(t *testing.T) {
				_, err := ValidateToken(tc.token, keys.AccessVerify)

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
