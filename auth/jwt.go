package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// TokenDetails represents tokens details needed for authentication
type TokenDetails struct {
	Access           string `json:"access_token"`
	Refresh          string `json:"refresh_token"`
	AccessUUID       string `json:"-"`
	RefreshUUID      string `json:"-"`
	AccessExpiresAt  int64  `json:"-"`
	RefreshExpiresAt int64  `json:"-"`
}

// Token creates access and refresh tokens for a user with specified ID
func Token(userID int) (*TokenDetails, error) {
	var err error
	details := &TokenDetails{}

	details.AccessExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	details.RefreshExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	details.AccessUUID = uuid.New().String()
	details.RefreshUUID = uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     details.AccessExpiresAt,
		"uuid":    details.AccessUUID,
	})
	details.Access, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     details.RefreshExpiresAt,
		"uuid":    details.RefreshUUID,
	})
	details.Refresh, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return details, nil
}
