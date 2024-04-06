package auth

import (
	"testing"
	"time"

	"github.com/mbeka02/bank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	username := utils.RandName()

	maker, err := NewPasetoMaker(utils.RandString(32))
	require.NoError(t, err)
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.ValidateToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)
	require.NotZero(t, claims.ID)
	require.Equal(t, username, claims.Username)
	require.WithinDuration(t, issuedAt, claims.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, claims.ExpiresAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandString(32))
	require.NoError(t, err)
	username := utils.RandName()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.ValidateToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, claims)
}
