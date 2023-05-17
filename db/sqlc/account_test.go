package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/takipipo/simple-bank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)

}
func TestGetAccount(t *testing.T) {
	expectedAccount := createRandomAccount(t)

	actualAccount, err := testQueries.GetAccount(context.Background(), expectedAccount.ID)

	require.NoError(t, err)
	require.Equal(t, expectedAccount.ID, actualAccount.ID)
	require.Equal(t, expectedAccount.Owner, actualAccount.Owner)
	require.Equal(t, expectedAccount.Currency, actualAccount.Currency)
	require.Equal(t, expectedAccount.CreatedAt, actualAccount.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomMoney(),
	}
	actualUpdatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, createdAccount.ID, actualUpdatedAccount.ID)
	require.Equal(t, createdAccount.Owner, actualUpdatedAccount.Owner)
	require.Equal(t, createdAccount.CreatedAt, actualUpdatedAccount.CreatedAt)
	require.Equal(t, createdAccount.Currency, actualUpdatedAccount.Currency)
	require.Equal(t, arg.Balance, actualUpdatedAccount.Balance)

}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)

	actualAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, actualAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	actualAccounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, actualAccounts, 5)

}
