package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// an interface that manages tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	ValidateToken(tokenString string) (jwt.MapClaims, error)
}
