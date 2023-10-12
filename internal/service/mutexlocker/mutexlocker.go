package mutexlocker

import (
	"context"
	"sync"
)

type L struct {
	// TODO: such "local-state-only" mutex doesn't make any sense
	// It is ok for testing purpose, but in future we should use some global state mutex
	// so multiple instance of L could have access to the same mutexes in memory
	m map[int]sync.Locker

	// in order to protect our map
	lock sync.Locker
}

func New() *L {
	return &L{
		m:    make(map[int]sync.Locker, 0),
		lock: &sync.Mutex{},
	}
}

func (l *L) getLock(userID int) sync.Locker {
	l.lock.Lock()
	defer l.lock.Unlock()

	m, ok := l.m[userID]
	if !ok {
		l.m[userID] = &sync.Mutex{}
		return l.m[userID]
	}

	return m
}

func (l *L) Lock(_ context.Context, userID int) error {
	l.getLock(userID).Lock()

	return nil
}

func (l *L) Unlock(_ context.Context, userID int) error {
	l.getLock(userID).Unlock()

	return nil
}
