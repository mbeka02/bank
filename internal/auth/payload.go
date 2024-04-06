package auth

import (
	"time"

	"github.com/google/uuid"
)

//var ErrExpiredToken = errors.New("The token has expired")

// payload data of the token
type Payload struct {
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	ID        uuid.UUID `json:"id"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        uuid,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}, nil
}

// Checks if token payload is valid
/*func (payload *Payload) Valid() error {
	// check expiration time
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	}
	return nil
}*/
