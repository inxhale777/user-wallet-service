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
	TransactionStatusHold      TransactionStatus = "hold"
	TransactionStatusComplete  TransactionStatus = "complete"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// UserLocker - thing what can acquire or wait&acquire lock on some key, userID in our case
type UserLocker interface {
	Lock(ctx context.Context, userID int) error
	Unlock(ctx context.Context, userID int) error
}

type TransactionRepo interface {
	Get(ctx context.Context, transactionID int) (*Transaction, error)
	Create(ctx context.Context, userID int, amount int, status TransactionStatus) (transactionID int, e error)
	Change(ctx context.Context, transactionID int, status TransactionStatus) error
	Total(ctx context.Context, userID int) (balance int, e error)
}

type WalletService interface {
	Balance(ctx context.Context, userID int) (balance int, e error)
	Deposit(ctx context.Context, userID int, amount int) error
	Hold(ctx context.Context, userID int, amount int) error
	Charge(ctx context.Context, transactionID int) error
	Cancel(ctx context.Context, transactionID int) error
}
