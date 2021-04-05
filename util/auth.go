package util

import (
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken takes a jwt.Claims interface and returns its corresponding signed string
func GenerateToken(claims jwt.Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, tokenError := token.SignedString(secretKey)

	return tokenString, tokenError
}
