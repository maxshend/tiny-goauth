package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
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

// RSAKeys contains private and public keys
type RSAKeys struct {
	AccessSign    *rsa.PrivateKey
	AccessVerify  *rsa.PublicKey
	RefreshSign   *rsa.PrivateKey
	RefreshVerify *rsa.PublicKey
}

type authErr string

func (e authErr) Error() string { return string(e) }

const (
	errEmptySecret   = authErr("Token secret is empty")
	errEmptyPassword = authErr("Password is empty")
)

// Token creates access and refresh tokens for a user with specified ID
func Token(userID int64, roles []string, keys *RSAKeys) (*TokenDetails, error) {
	var err error

	details := &TokenDetails{}

	details.AccessExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	details.RefreshExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()

	details.AccessUUID = uuid.New().String()
	details.RefreshUUID = uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		userID,
		roles,
		details.AccessUUID,
		jwt.StandardClaims{
			ExpiresAt: details.AccessExpiresAt,
		},
	})
	details.Access, err = accessToken.SignedString(keys.AccessSign)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		userID,
		roles,
		details.RefreshUUID,
		jwt.StandardClaims{
			ExpiresAt: details.RefreshExpiresAt,
		},
	})
	details.Refresh, err = refreshToken.SignedString(keys.RefreshSign)
	if err != nil {
		return nil, err
	}

	return details, nil
}

// ValidateToken validates access and refresh tokens
func ValidateToken(tokenString string, secret *rsa.PublicKey) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		m, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok || m.Alg() != "RS256" {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
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

// Keys generates access and refresh RSA keys
func Keys() (*RSAKeys, error) {
	var err error
	keys := &RSAKeys{}

	bytes, err := ioutil.ReadFile(os.Getenv("ACCESS_PRIVATE_PATH"))
	if err != nil {
		return nil, err
	}

	keys.AccessSign, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = ioutil.ReadFile(os.Getenv("ACCESS_PUBLIC_PATH"))
	if err != nil {
		return nil, err
	}

	keys.AccessVerify, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = ioutil.ReadFile(os.Getenv("REFRESH_PRIVATE_PATH"))
	if err != nil {
		return nil, err
	}

	keys.RefreshSign, err = jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = ioutil.ReadFile(os.Getenv("REFRESH_PUBLIC_PATH"))
	if err != nil {
		return nil, err
	}

	keys.RefreshVerify, err = jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
