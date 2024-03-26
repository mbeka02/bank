package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := &Store{
		testQueries,
		testDB,
	}

	sender := createRandomAccount(t)
	receiver := createRandomAccount(t)
	/*the transactions are being ran on separate go-routines so channels are needed to send back the info to the main go-routine */
	testError := make(chan error)
	testResult := make(chan TransferTxResult)

	//run n concurrent transactions
	n := 4
	amount := int64(100)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				SenderID:   sender.ID,
				ReceiverID: receiver.ID,
				Amount:     amount,
			})
			//send  to channels
			testError <- err
			testResult <- result
		}()
	}
	// receive the stuff sent over the channels
	for i := 0; i < n; i++ {
		err := <-testError
		result := <-testResult
		require.NoError(t, err)
		require.NotEmpty(t, result)

		transfer := result.Transfer

		require.NotEmpty(t, transfer)

		require.Equal(t, sender.ID, transfer.SenderID)
		require.Equal(t, receiver.ID, transfer.ReceiverID)

		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		//confirm transfer record exists in the database
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		senderEntry := result.SenderEntry
		require.NotEmpty(t, senderEntry)
		require.Equal(t, sender.ID, senderEntry.AccountID)
		require.Equal(t, -amount, senderEntry.Amount)
		require.NotZero(t, senderEntry.ID)
		require.NotZero(t, senderEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), senderEntry.ID)
		require.NoError(t, err)

		receiverEntry := result.ReceiverEntry
		require.NotEmpty(t, receiverEntry)
		require.Equal(t, receiver.ID, receiverEntry.AccountID)
		require.Equal(t, amount, receiverEntry.Amount)
		require.NotZero(t, receiverEntry.ID)
		require.NotZero(t, receiverEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), receiverEntry.ID)
		require.NoError(t, err)

		//TO DO : CHECK ACCOUNT BALANCES
	}

}
