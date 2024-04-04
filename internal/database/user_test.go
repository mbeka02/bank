package database

import (
	"context"
	"testing"
	"time"

	"github.com/mbeka02/bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	args := CreateUserParams{
		UserName: utils.RandName(),
		FullName: utils.RandName(),
		Password: "password",
		Email:    utils.RandEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user

}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), randomUser.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, randomUser.UserName, user.UserName)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.Password, user.Password)
	require.Equal(t, randomUser.Email, user.Email)
	require.WithinDuration(t, randomUser.CreatedAt, user.CreatedAt, time.Second)

}
