package domain

import "context"

type Transaction struct {
	ID        int
	UserID    int
	ServiceID int
	OrderID   int
	Status    TransactionStatus
	Amount    int
}

type TransactionRepo interface {
	Get(ctx context.Context, transactionID int) (*Transaction, error)
	Create(ctx context.Context, userID int, amount int, status TransactionStatus) (transactionID int, e error)
	Change(ctx context.Context, transactionID int, status TransactionStatus) error
	Total(ctx context.Context, userID int) (balance int, e error)
}
