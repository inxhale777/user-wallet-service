package pgtransactions

import (
	"context"
	"fmt"
	"user-wallet-service/internal/domain"
	"user-wallet-service/internal/postgres"

	"github.com/jackc/pgx/v5/pgtype"
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
		QueryRow(ctx, "select id, user_id, status, amount, description from transactions where id = $1", transactionID).
		Scan(&tx.ID, &tx.UserID, &tx.Status, &tx.Amount, &tx.Description)
	if err != nil {
		return nil, fmt.Errorf("pg_transactions.Get: %w", err)
	}

	if tx.ID == 0 {
		return nil, fmt.Errorf("pg_transactions.Get: %w", domain.NewErrTxNotFound(transactionID))
	}

	return &tx, nil
}

func (t *R) Create(ctx context.Context, tx domain.Transaction) (int, error) {
	var id int
	err := t.db.QueryRow(ctx,
		"insert into transactions (user_id, amount, status, description) values ($1, $2, $3, $4) returning id;",
		tx.UserID, tx.Amount, tx.Status, tx.Description).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("pg_transactions.Create: %w", err)
	}

	return id, nil
}
func (t *R) Change(ctx context.Context, transactionID int, status domain.TransactionStatus) error {
	sql := "update transactions set status = $1 where id = $2 and status = $3"
	r, err := t.db.Exec(ctx, sql, status, transactionID, domain.TransactionStateMachine[status])
	if err != nil {
		return fmt.Errorf("pg_transactions.Change: %w", err)
	}

	if r.RowsAffected() == 0 {
		return fmt.Errorf("pg_transactions.Change: %w", domain.NewErrTxNotFound(transactionID))
	}

	return nil
}

func (t *R) Total(ctx context.Context, userID int) (int, error) {
	var total pgtype.Int8
	err := t.db.QueryRow(ctx,
		"select sum(amount) from transactions where user_id = $1 and status in ($2, $3)",
		userID, domain.TransactionStatusComplete, domain.TransactionStatusHold).
		Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("pg_transactions.Total: %w", err)
	}

	return int(total.Int64), nil
}
