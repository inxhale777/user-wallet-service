package domain

import (
	"context"
)

type Transaction struct {
	ID        int
	UserID    int
	ServiceID int
	OrderID   int
	Status    TransactionStatus
	Amount    int
}

type TransactionStatus string

var (
	TransactionStatusHold     TransactionStatus = "hold"
	TransactionStatusComplete TransactionStatus = "complete"
	TransactionStatusCanceled TransactionStatus = "canceled"
)

// UserLocker - thing what can acquire or wait&acquire lock on some key, userID in our case
type UserLocker interface {
	Lock(ctx context.Context, userID string) error
}

type TransactionRepo interface {
	Get(ctx context.Context, transactionID string) (*Transaction, error)
	Create(ctx context.Context, userID string, amount int, status TransactionStatus) (transactionID string, e error)
	Change(ctx context.Context, transactionID string, status TransactionStatus) error
	Total(ctx context.Context, userID string) (balance int, e error)
}

type WalletService interface {
	Balance(ctx context.Context, userID string) (balance int, e error)
	Deposit(ctx context.Context, userID string, amount int) error
	Hold(ctx context.Context, userID string, amount int) error
	Charge(ctx context.Context, transactionID string) error
	Cancel(ctx context.Context, transactionID string) error
}
