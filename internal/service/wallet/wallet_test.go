package wallet_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"user-balance-service/internal/mock/mock_transactions"
	"user-balance-service/internal/service/wallet"
)

func TestWallet_DepositWithoutLocker(t *testing.T) {
	userID := 9999
	amount := 500 * 100

	w := wallet.New(mock_transactions.New(), nil)
	err := w.Deposit(context.Background(), userID, amount)
	require.ErrorContains(t, err, "you must provide locker impl. to use this method")
}
