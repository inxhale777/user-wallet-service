package domain

import "context"

// UserLocker - thing what can acquire or wait&acquire lock on some key, userID in our case
type UserLocker interface {
	Lock(ctx context.Context, userID int) error
	Unlock(ctx context.Context, userID int) error
}
