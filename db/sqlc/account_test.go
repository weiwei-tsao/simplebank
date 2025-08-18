package db

import (
	"context"
	"simplebank/db/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  pgtype.Int8{Int64: util.RandomMoney(), Valid: true},
		Currency: pgtype.Text{String: util.RandomCurrency(), Valid: true},
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// check that no error occurred
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// check that the account is created with the correct values
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
