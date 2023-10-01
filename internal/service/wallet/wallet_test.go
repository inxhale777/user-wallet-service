package wallet_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"user-balance-service/internal/repo/inmemory_transactions"
	"user-balance-service/internal/service/mutex_locker"
	"user-balance-service/internal/service/wallet"
)

func TestWallet_Deposit(t *testing.T) {

	t.Run("without locker", func(t *testing.T) {
		ctx := context.Background()
		userID := 9999
		amount := 500 * 100

		w := wallet.New(inmemory_transactions.New(), nil)
		err := w.Deposit(ctx, userID, amount)
		require.ErrorContains(t, err, "you must provide locker impl. to use this method")
	})

	t.Run("concurrently", func(t *testing.T) {
		ctx := context.Background()
		runs := 999
		userID := 9999
		amount := 500 * 100

		w := wallet.New(inmemory_transactions.New(), mutex_locker.New())

		// actually we can't get any anomalies doing .Deposit concurrently
		// but let it be for fun
		var wg sync.WaitGroup
		for i := 0; i < runs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := w.Deposit(ctx, userID, amount)
				require.Nil(t, err)
			}()
		}

		wg.Wait()

		balance, err := w.Balance(ctx, userID)
		require.Nil(t, err)
		require.Equal(t, amount*runs, balance)
	})
}

func TestWallet_Hold(t *testing.T) {

	t.Run("without locker", func(t *testing.T) {
		ctx := context.Background()
		userID := 9999
		amount := 500 * 100

		w := wallet.New(inmemory_transactions.New(), nil)
		err := w.Hold(ctx, userID, amount)
		require.ErrorContains(t, err, "you must provide locker impl. to use this method")
	})

	t.Run("concurrently", func(t *testing.T) {
		ctx := context.Background()
		userID := 9999
		hold := 500 * 100
		// 50.000 $
		balance := hold * 100
		// 100 of holds must be success
		// but on another 200 we don't have enough money
		runs := 100 + 200

		w := wallet.New(inmemory_transactions.New(), mutex_locker.New())
		err := w.Deposit(ctx, userID, balance)
		require.Nil(t, err)

		var success atomic.Int32
		var failed atomic.Int32
		var wg sync.WaitGroup
		for i := 0; i < runs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := w.Hold(ctx, userID, hold)
				if err != nil {
					require.ErrorContains(t, err, "not enough money, have")
					failed.Add(1)
					return
				}

				success.Add(1)
			}()
		}
		wg.Wait()

		b, err := w.Balance(ctx, userID)
		require.Nil(t, err)
		require.Equal(t, 0, b)

		require.Equal(t, int32(100), success.Load())
		require.Equal(t, int32(200), failed.Load())
	})
}
