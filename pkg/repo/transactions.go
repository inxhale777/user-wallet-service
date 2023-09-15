package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"user-balance-service/pkg/domain"
	"user-balance-service/pkg/postgres"
)

type TransactionRepo struct {
	db postgres.DB
}

func NewTransactionPGRepo(db postgres.DB) *TransactionRepo {
	return &TransactionRepo{db}
}

func (t *TransactionRepo) Get(ctx context.Context, transactionID string) (*domain.Transaction, error) {

	var tx domain.Transaction
	err := t.db.
		QueryRow(ctx, "select id, user_id, service_id, order_id, status, amount from transactions where id = $1", transactionID).
		Scan(&tx.ID, &tx.UserID, &tx.ServiceID, &tx.OrderID, &tx.Status, &tx.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "TransactionRepo.Get")
	}

	return &tx, nil
}

func (t *TransactionRepo) Create(ctx context.Context, userID string, amount int, status domain.TransactionStatus) (string, error) {
	var id string
	err := t.db.QueryRow(ctx,
		"insert into transactions (user_id, amount, status) values ($1, $2, $3) returning id;", userID, amount, status).
		Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, "TransactionRepo.Create")
	}

	return id, nil
}

func (t *TransactionRepo) Change(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
	_, err := t.db.Exec(ctx, "update transactions set status = $1 where id = $2", status, transactionID)
	if err != nil {
		return errors.Wrap(err, "TransactionRepo.Change")
	}

	return nil
}

func (t *TransactionRepo) Total(ctx context.Context, userID string) (balance int, e error) {
	var total pgtype.Int8
	err := t.db.QueryRow(ctx,
		"select sum(amount) from transactions where user_id = $1 and status in ($2, $3)",
		userID, domain.TransactionStatusComplete, domain.TransactionStatusHold).
		Scan(&total)
	if err != nil {
		return 0, errors.Wrap(err, "TransactionRepo.Total")
	}

	return int(total.Int64), nil
}
