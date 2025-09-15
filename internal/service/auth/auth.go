package authService

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTAuthenticatorService struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticatorService(secret, aud, iss string) *JWTAuthenticatorService {
	return &JWTAuthenticatorService{secret: secret, aud: aud, iss: iss}
}

func (a *JWTAuthenticatorService) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (a *JWTAuthenticatorService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	})
}
