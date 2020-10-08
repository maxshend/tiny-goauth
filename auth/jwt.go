package auth

import (
	"fmt"
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

// Claims represents data from JWT body
type Claims struct {
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
	UUID   string   `json:"uuid"`
	jwt.StandardClaims
}

type authErr string

func (e authErr) Error() string { return string(e) }

const (
	errEmptySecret   = authErr("Token secret is empty")
	errEmptyPassword = authErr("Password is empty")
)

// Token creates access and refresh tokens for a user with specified ID
func Token(userID int64, roles []string) (*TokenDetails, error) {
	var err error

	if len(os.Getenv("ACCESS_TOKEN_SECRET")) == 0 || len(os.Getenv("REFRESH_TOKEN_SECRET")) == 0 {
		return nil, errEmptySecret
	}

	details := &TokenDetails{}

	details.AccessExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	details.RefreshExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	details.AccessUUID = uuid.New().String()
	details.RefreshUUID = uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		userID,
		roles,
		details.AccessUUID,
		jwt.StandardClaims{
			ExpiresAt: details.AccessExpiresAt,
		},
	})
	details.Access, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		userID,
		roles,
		details.RefreshUUID,
		jwt.StandardClaims{
			ExpiresAt: details.RefreshExpiresAt,
		},
	})
	details.Refresh, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return details, nil
}

// ValidateAccessToken validates access token
func ValidateAccessToken(tokenString string) (jwt.Claims, error) {
	return validateToken(tokenString, os.Getenv("ACCESS_TOKEN_SECRET"))
}

// ValidateRefreshToken validates refresh token
func ValidateRefreshToken(tokenString string) (jwt.Claims, error) {
	return validateToken(tokenString, os.Getenv("REFRESH_TOKEN_SECRET"))
}

func validateToken(tokenString, secret string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		hmac, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || hmac.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
