package service

import (
	"context"
	"github.com/pkg/errors"
	"user-balance-service/pkg/domain"
)

type WalletService struct {
	txRepo domain.TransactionRepo
}

func NewWallet(txRepo domain.TransactionRepo) *WalletService {
	return &WalletService{txRepo}
}

func (w *WalletService) Balance(ctx context.Context, userID string) (balance int, e error) {
	b, err := w.txRepo.Total(ctx, userID)
	if err != nil {
		return 0, errors.Wrap(err, "WalletService.Balance")
	}

	return b, nil
}

func (w *WalletService) Deposit(ctx context.Context, userID string, amount int) error {
	_, err := w.txRepo.Create(ctx, userID, amount, domain.TransactionStatusComplete)
	if err != nil {
		return errors.Wrap(err, "WalletService.Deposit")
	}

	return nil
}

func (w *WalletService) Hold(ctx context.Context, userID string, amount int) error {
	//TODO implement me
	panic("implement me")
}

func (w *WalletService) Charge(ctx context.Context, transactionID string) error {
	//TODO implement me
	panic("implement me")
}

func (w *WalletService) Cancel(ctx context.Context, transactionID string) error {
	//TODO implement me
	panic("implement me")
}
