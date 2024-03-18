package db

import (
	"context"
	"testing"

	"github.com/dsaladbar617/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount: util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, entry.Amount, arg.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {

	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.ID, entry2.ID)

}

func createTestAccount(t *testing.T) Account {
	accountArg := CreateAccountParams{
		Owner: "Dan",
		Balance: 1000,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), accountArg)
	require.NoError(t, err)

	return account
}

func createTestEntry(t *testing.T, account Account) {

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: 1000,
	}

	_, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	// return entry
}

func TestListEntries(t *testing.T) {

	account1 := createTestAccount(t)


	for i := 0; i < 10; i++ {
		createTestEntry(t, account1)
	}

	arg := ListEntriesParams{
		AccountID: account1.ID,
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}