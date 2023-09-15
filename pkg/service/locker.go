package service

import (
	"context"
	"user-balance-service/pkg/postgres"
)

type Locker struct {
	// can be used ONLY inside pgsql transaction
	tx postgres.Tx
}

func NewLocker(tx postgres.Tx) *Locker {
	return &Locker{tx}
}

func (l *Locker) Lock(ctx context.Context, userID string) error {
	_, err := l.tx.Exec(ctx, "select pg_advisory_xact_lock($1);", userID)
	if err != nil {
		return err
	}

	return nil
}
