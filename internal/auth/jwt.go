package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretLength = 32

var ErrInvalidToken = errors.New("token is invalid")
var ErrExpiredToken = errors.New("the token has expired")

type JWTMaker struct {
	secret string
}

// ensure JWTMaker implements Maker
func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < minSecretLength {
		return nil, errors.New("secret is too short , it must be atleast 32 characters long")
	}
	return &JWTMaker{secret}, nil

}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secret))
}

func (maker *JWTMaker) ValidateToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(maker.secret), nil
	} /* Parse takes the token string and a function for looking up the key. The latter is especially useful if you use multiple keys for your application.  The standard is to use 'kid' in the head of the token to identify which key to use, but the parsed token (head and claims) is provided to the callback, providing flexibility.*/
	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := jwtToken.Claims.(*Payload)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, ErrExpiredToken
	}
	return claims, nil
}
