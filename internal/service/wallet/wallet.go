package wallet

import (
	"context"
	"fmt"
	"user-wallet-service/internal/domain"
)

type W struct {
	txRepo domain.TransactionRepo
	locker domain.UserLocker
}

func New(txRepo domain.TransactionRepo, locker domain.UserLocker) *W {
	return &W{txRepo, locker}
}

func (w *W) Balance(ctx context.Context, userID int) (int, error) {
	b, err := w.txRepo.Total(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("wallet.Balance: %w", domain.ErrInvalidAmount)
	}

	return b, nil
}

func (w *W) Deposit(ctx context.Context, userID int, amount int) error {

	trace := "wallet.Deposit"

	if w.locker == nil {
		return fmt.Errorf("%s: %w", trace, domain.ErrNoLockerProvided)
	}

	err := w.locker.Lock(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", trace, err)
	}
	defer func() {
		_ = w.locker.Unlock(ctx, userID)
	}()

	_, err = w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusComplete)
	if err != nil {
		return fmt.Errorf("%s: %w", trace, err)
	}

	return nil
}

func (w *W) Hold(ctx context.Context, userID int, amount int) (int, error) {

	trace := "wallet.Hold"

	if w.locker == nil {
		return 0, fmt.Errorf("%s: %w", trace, domain.ErrNoLockerProvided)
	}

	if amount < 1 {
		return 0, fmt.Errorf("%s: %w", trace, domain.ErrInvalidAmount)
	}

	err := w.locker.Lock(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", trace, err)
	}
	defer func() {
		_ = w.locker.Unlock(ctx, userID)
	}()

	balance, err := w.txRepo.Total(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", trace, err)
	}

	if balance < amount {
		return 0, fmt.Errorf("%s: %w", trace, domain.NewErrInsufficientMoney(userID, amount, balance))
	}

	// create transaction with NEGATIVE value because we are CHARGING money from user
	amount *= -1
	tx, err := w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusHold)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", trace, err)
	}

	return tx, nil
}

func (w *W) Charge(ctx context.Context, transactionID int) error {
	trace := "wallet.Charge"

	tx, err := w.txRepo.Get(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("%s: %w", trace, err)
	}

	want := domain.TransactionStatusComplete
	if tx.Status != domain.TransactionStatusHold {
		return fmt.Errorf("%s: %w", trace, domain.NewErrInvalidTxStatus(transactionID, want, tx.Status))
	}

	// do we need locker here?
	err = w.txRepo.Change(ctx, transactionID, want)
	if err != nil {
		switch err.(type) {
		case *domain.ErrTxNotFound:
			// Someone has changed TX faster than us.
			// That's why we got ErrTxNotFound despite of fact that we have found it before
			// So grab TX again with newer status and return ErrInvalidTxStatus
			tx, err = w.txRepo.Get(ctx, transactionID)
			if err != nil {
				return fmt.Errorf("%s: %w", trace, err)
			}

			return fmt.Errorf("%s: %w", trace, domain.NewErrInvalidTxStatus(transactionID, want, tx.Status))
		default:
			return fmt.Errorf("%s: %w", trace, err)
		}
	}

	return nil
}

func (w *W) Cancel(ctx context.Context, transactionID int) error {
	trace := "wallet.Cancel"

	tx, err := w.txRepo.Get(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("%s: %w", trace, err)
	}

	want := domain.TransactionStatusCancelled
	if tx.Status != domain.TransactionStatusHold {
		return fmt.Errorf("%s: %w", trace, domain.NewErrInvalidTxStatus(transactionID, want, tx.Status))
	}

	// do we need locker here?
	err = w.txRepo.Change(ctx, transactionID, want)
	if err != nil {
		return fmt.Errorf("%s: %w", trace, err)
	}

	return nil
}
