package domain

import "context"

type Transaction struct {
	ID          int
	UserID      int
	Amount      int
	Status      TransactionStatus
	Description string
}

type TransactionRepo interface {
	Get(ctx context.Context, transactionID int) (*Transaction, error)
	Create(ctx context.Context, transaction Transaction) (transactionID int, e error)
	Change(ctx context.Context, transactionID int, status TransactionStatus) error
	Total(ctx context.Context, userID int) (balance int, e error)
}
