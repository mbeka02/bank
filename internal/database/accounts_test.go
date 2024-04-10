package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mbeka02/bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:   user.UserName,
		Balance: utils.RandMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)

}

func TestGetAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, randomAccount.ID, account.ID)
	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Currency, account.Currency)
	require.Equal(t, randomAccount.Balance, account.Balance)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	newBalance := utils.RandMoney()
	account, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      randomAccount.ID,
		Balance: newBalance,
	})
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, randomAccount.ID, account.ID)
	require.Equal(t, randomAccount.Owner, account.Owner)
	require.Equal(t, randomAccount.Currency, account.Currency)
	require.Equal(t, newBalance, account.Balance)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), randomAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestGetAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	args := GetAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.GetAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
