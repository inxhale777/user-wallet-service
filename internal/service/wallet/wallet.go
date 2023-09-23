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

	_, err = w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusComplete)
	if err != nil {
		return errors.Wrap(err, trace)
	}

	err = w.locker.Unlock(ctx, userID)
	if err != nil {
		return errors.Wrap(err, trace)
	}

	return nil
}

func (w *W) Hold(ctx context.Context, userID int, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (w *W) Charge(ctx context.Context, transactionID int) error {
	//TODO implement me
	panic("implement me")
}

func (w *W) Cancel(ctx context.Context, transactionID int) error {
	//TODO implement me
	panic("implement me")
}
