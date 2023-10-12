package pglocker

import (
	"context"
	"fmt"
	"user-wallet-service/internal/postgres"
)

type Tx interface {
	postgres.DB
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Locker struct {
	// can be used ONLY inside pgsql transaction
	tx Tx
}

func New(tx Tx) *Locker {
	return &Locker{tx}
}

func (l *Locker) Lock(ctx context.Context, userID int) error {
	_, err := l.tx.Exec(ctx, "select pg_advisory_xact_lock($1);", userID)
	if err != nil {
		return fmt.Errorf("pg_locker.Lock: %w", err)
	}

	return nil
}

func (l *Locker) Unlock(context.Context, int) error {
	// we use pg_advisory_xact_lock here
	// so it will automatically unlocked
	// after commit or rollback of postgresql transaction

	return nil
}
