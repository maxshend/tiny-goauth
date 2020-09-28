package auth

import (
	"fmt"
	"log"
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
	UserID int64  `json:"user_id"`
	UUID   string `json:"uuid"`
	jwt.StandardClaims
}

// Token creates access and refresh tokens for a user with specified ID
func Token(userID int64) (*TokenDetails, error) {
	var err error
	details := &TokenDetails{}

	details.AccessExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	details.RefreshExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	details.AccessUUID = uuid.New().String()
	details.RefreshUUID = uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		userID,
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
		details.AccessUUID,
		jwt.StandardClaims{
			ExpiresAt: details.AccessExpiresAt,
		},
	})
	details.Refresh, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return details, nil
}

// ValidateToken validates access token
func ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		hmac, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || hmac.Alg() != "HS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	log.Println(claims.UserID)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
