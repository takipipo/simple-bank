package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// Concurrent one-directional transfer account1 -> account2, account2 -> account1
func TestTransferTxOneDirectional(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        amount,
				},
			)

			errs <- err // channel <- value-in-different-go-routine
			results <- result
		}()
	}
	// check results
	for i := 0; i < n; i++ {
		err := <-errs // variable-in-main-go-routine <- channel
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)

		// check accounts
		actualFromAccount := result.FromAccount
		require.NotEmpty(t, actualFromAccount)
		require.Equal(t, fromAccount.ID, actualFromAccount.ID)

		actualToAccount := result.ToAccount
		require.NotEmpty(t, actualToAccount)
		require.Equal(t, toAccount.ID, actualToAccount.ID)

		// check accounts' balance
		diff1 := fromAccount.Balance - actualFromAccount.Balance
		diff2 := actualToAccount.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		k := int(diff1 / amount)
		require.True(t, 1 <= k && k <= n)
	}
	// check the final updated balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)

}

// Concurrent bi-directional transfer account1 -> account2, account2 -> account1
func TestTransferTxBiDirectional(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// half of the concurrent have to transfer in another direction
		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(
				context.Background(),
				TransferTxParams{
					FromAccountID: fromAccountID,
					ToAccountID:   toAccountID,
					Amount:        amount,
				},
			)

			errs <- err // channel <- value-in-different-go-routine
		}()
	}
	// check results
	for i := 0; i < n; i++ {
		err := <-errs // variable-in-main-go-routine <- channel
		require.NoError(t, err)
	}
	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
