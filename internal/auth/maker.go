package auth

import (
	"time"
)

// an interface that manages tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	ValidateToken(tokenString string) (*Payload, error)
}
