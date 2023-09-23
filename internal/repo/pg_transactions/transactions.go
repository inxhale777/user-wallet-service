package pg_transactions

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"user-balance-service/internal/domain"
	"user-balance-service/internal/postgres"
)

type R struct {
	db postgres.DB
}

func New(db postgres.DB) *R {
	return &R{db}
}

func (t *R) Get(ctx context.Context, transactionID int) (*domain.Transaction, error) {

	var tx domain.Transaction
	err := t.db.
		QueryRow(ctx, "select id, user_id, service_id, order_id, status, amount from transactions where id = $1", transactionID).
		Scan(&tx.ID, &tx.UserID, &tx.ServiceID, &tx.OrderID, &tx.Status, &tx.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "R.Get")
	}

	return &tx, nil
}

func (t *R) Create(ctx context.Context, userID int, amount int, status domain.TransactionStatus) (int, error) {
	var id int
	err := t.db.QueryRow(ctx,
		"insert into transactions (user_id, amount, status) values ($1, $2, $3) returning id;", userID, amount, status).
		Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "R.Create")
	}

	return id, nil
}

func (t *R) Change(ctx context.Context, transactionID int, status domain.TransactionStatus) error {
	_, err := t.db.Exec(ctx, "update transactions set status = $1 where id = $2", status, transactionID)
	if err != nil {
		return errors.Wrap(err, "R.Change")
	}

	return nil
}

func (t *R) Total(ctx context.Context, userID int) (balance int, e error) {
	var total pgtype.Int8
	err := t.db.QueryRow(ctx,
		"select sum(amount) from transactions where user_id = $1 and status in ($2, $3)",
		userID, domain.TransactionStatusComplete, domain.TransactionStatusHold).
		Scan(&total)
	if err != nil {
		return 0, errors.Wrap(err, "R.Total")
	}

	return int(total.Int64), nil
}
