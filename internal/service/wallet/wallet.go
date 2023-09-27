package wallet

import (
	"context"
	"github.com/pkg/errors"
	"user-balance-service/internal/domain"
)

type W struct {
	txRepo domain.TransactionRepo
	locker domain.UserLocker
}

func New(txRepo domain.TransactionRepo, locker domain.UserLocker) *W {
	return &W{txRepo, locker}
}

func (w *W) Balance(ctx context.Context, userID int) (balance int, e error) {
	b, err := w.txRepo.Total(ctx, userID)
	if err != nil {
		return 0, errors.Wrap(err, "WalletService.Balance")
	}

	return b, nil
}

func (w *W) Deposit(ctx context.Context, userID int, amount int) error {

	trace := "WalletService.Deposit"

	if w.locker == nil {
		return errors.New("you must provide locker impl. to use this method")
	}

	err := w.locker.Lock(ctx, userID)
	if err != nil {
		return errors.Wrap(err, trace)
	}
	defer func() {
		_ = w.locker.Unlock(ctx, userID)
	}()

	_, err = w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusComplete)
	if err != nil {
		return errors.Wrap(err, trace)
	}

	return nil
}

func (w *W) Hold(ctx context.Context, userID int, amount int) error {

	trace := "WalletService.Hold"

	if w.locker == nil {
		return errors.New("you must provide locker impl. to use this method")
	}

	err := w.locker.Lock(ctx, userID)
	if err != nil {
		return errors.Wrap(err, trace)
	}
	defer func() {
		_ = w.locker.Unlock(ctx, userID)
	}()

	balance, err := w.txRepo.Total(ctx, userID)
	if err != nil {
		return errors.Wrap(err, trace)
	}

	if balance < amount {
		return errors.Errorf("not enough money, have %d, but want to hold: %d", balance/100, amount/100)
	}

	// create transaction with NEGATIVE value because we are CHARGING money from user
	amount *= -1
	_, err = w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusHold)
	if err != nil {
		return errors.Wrap(err, trace)
	}

	return nil
}

func (w *W) Charge(ctx context.Context, transactionID int) error {
	//TODO implement me
	panic("implement me")
}

func (w *W) Cancel(ctx context.Context, transactionID int) error {
	//TODO implement me
	panic("implement me")
}
