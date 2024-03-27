package database

import (
	"context"
	"log"
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

	//log the balances before sending
	log.Printf("The current account balance for the sender is : %v \n", sender.Balance)
	log.Printf("The current account balance for the receiver is : %v \n", receiver.Balance)
	/*the transactions are being ran on separate go-routines so channels are needed to send back the info to the main go-routine */
	testError := make(chan error)
	testResult := make(chan TransferTxResult)

	//run n concurrent transactions
	n := 5
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
	hasExisted := make(map[int]bool)
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

		//check accounts
		senderAccount := result.SenderAccount
		require.NotEmpty(t, senderAccount)
		require.Equal(t, sender.ID, senderAccount.ID)

		receiverAccount := result.ReceiverAccount
		require.NotEmpty(t, receiverAccount)
		require.Equal(t, receiver.ID, receiverAccount.ID)

		// CHECK ACCOUNT BALANCES
		diff1 := sender.Balance - senderAccount.Balance
		diff2 := receiverAccount.Balance - receiver.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //1*amount , 2*amount ... n*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		/*k also has to be unique for each transaction 1 , 2 , 3 ,4*/
		require.NotContains(t, hasExisted, k)
		hasExisted[k] = true
	}
	//check the final updated Balance
	UpdatedSenderAccount, err := testQueries.GetAccount(context.Background(), sender.ID)
	require.NoError(t, err)

	UpdatedReceiverAccount, err := testQueries.GetAccount(context.Background(), receiver.ID)
	//log the account balances after the transaction
	log.Printf("The account balance of the sender after sending is:%v \n", UpdatedSenderAccount.Balance)
	log.Printf("The account balance for the receiver after sending is:%v \n", UpdatedReceiverAccount.Balance)
	require.NoError(t, err)
	require.Equal(t, sender.Balance-int64(n)*amount, UpdatedSenderAccount.Balance)
	require.Equal(t, receiver.Balance+int64(n)*amount, UpdatedReceiverAccount.Balance)

}
