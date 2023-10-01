package wallet_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
	"user-wallet-service/internal/domain"
	"user-wallet-service/internal/repo/inmemory_transactions"
	"user-wallet-service/internal/service/mutex_locker"
	"user-wallet-service/internal/service/wallet"
)

func TestWallet_Deposit(t *testing.T) {

	t.Run("without locker", func(t *testing.T) {
		w := wallet.New(inmemory_transactions.New(), nil)
		err := w.Deposit(context.Background(), 0, 0)
		require.ErrorIs(t, err, domain.ErrNoLockerProvided)
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
		w := wallet.New(inmemory_transactions.New(), nil)
		_, err := w.Hold(context.Background(), 0, 0)
		require.ErrorIs(t, err, domain.ErrNoLockerProvided)
	})

	t.Run("concurrently", func(t *testing.T) {
		ctx := context.Background()
		userID := 9999
		hold := 500 * 100
		// 50.000 $
		balance := hold * 100
		// 100 of holds must be success
		// but on another 200 we don't have enough money
		success := int32(100)
		failed := int32(200)
		runs := int(success + failed)

		w := wallet.New(inmemory_transactions.New(), mutex_locker.New())
		err := w.Deposit(ctx, userID, balance)
		require.Nil(t, err)

		var s atomic.Int32
		var f atomic.Int32
		var wg sync.WaitGroup
		for i := 0; i < runs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := w.Hold(ctx, userID, hold)
				if err != nil {

					var insufficientMoneyErr *domain.ErrInsufficientMoney
					require.ErrorAs(t, err, &insufficientMoneyErr)
					require.Equal(t, insufficientMoneyErr.UserID, userID)
					require.Equal(t, insufficientMoneyErr.Want, hold)
					require.Equal(t, insufficientMoneyErr.Have, 0)

					f.Add(1)
					return
				}

				s.Add(1)
			}()
		}
		wg.Wait()

		b, err := w.Balance(ctx, userID)
		require.Nil(t, err)
		require.Equal(t, 0, b)

		require.Equal(t, success, s.Load())
		require.Equal(t, failed, f.Load())
	})
}

func TestWallet_Charge(t *testing.T) {

	t.Run("concurrently", func(t *testing.T) {

		ctx := context.Background()
		userID := 9999
		// 1000 $
		balance := 1000 * 100
		hold := balance
		success := int32(1)
		failed := int32(99)
		runs := int(success + failed)

		w := wallet.New(inmemory_transactions.New(), mutex_locker.New())
		err := w.Deposit(ctx, userID, balance)
		require.Nil(t, err)

		tx, err := w.Hold(ctx, userID, hold)
		require.Nil(t, err)

		var s atomic.Int32
		var f atomic.Int32
		var wg sync.WaitGroup
		for i := 0; i < runs; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := w.Charge(ctx, tx)
				if err != nil {

					var txStatusErr *domain.ErrInvalidTxStatus
					require.ErrorAs(t, err, &txStatusErr)
					require.Equal(t, txStatusErr.TransactionID, tx)

					// attempt to change status: complete ----> complete
					require.Equal(t, txStatusErr.Want, domain.TransactionStatusComplete)
					require.Equal(t, txStatusErr.Have, domain.TransactionStatusComplete)

					f.Add(1)
					return
				}

				s.Add(1)
			}()
		}
		wg.Wait()

		b, err := w.Balance(ctx, userID)
		require.Nil(t, err)
		require.Equal(t, 0, b)

		require.Equal(t, success, s.Load())
		require.Equal(t, failed, f.Load())
	})
}
