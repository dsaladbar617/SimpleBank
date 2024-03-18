package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestAccounts(t *testing.T) (Account, Account) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	return account1, account2
}

func createRandomTransfer(t *testing.T) Transfer {
	account1, account2 := createTestAccounts(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: 50,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)

	require.Equal(t, transfer.FromAccountID, account1.ID)
	require.Equal(t, transfer.ToAccountID, account2.ID)
	require.Equal(t, transfer.Amount, arg.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T){

	testTransfer := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), testTransfer.ID)

	require.NoError(t, err)
	require.Equal(t, testTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, testTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, testTransfer.Amount, transfer.Amount)
}

func createTestTransfer(t *testing.T, account1 Account, account2 Account) {

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: 100,
	}

	_, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
}

func TestListTransfers(t *testing.T) {
	account1, account2 := createTestAccounts(t)

	for i := 0; i < 10; i++ {
		createTestTransfer(t, account1, account2)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}