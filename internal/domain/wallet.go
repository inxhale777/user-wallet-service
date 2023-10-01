package domain

import "context"

type WalletService interface {
	Balance(ctx context.Context, userID int) (balance int, e error)
	Deposit(ctx context.Context, userID int, amount int) error
	Hold(ctx context.Context, userID int, amount int) (transactionID int, e error)
	Charge(ctx context.Context, transactionID int) error
	Cancel(ctx context.Context, transactionID int) error
}
