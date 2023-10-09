package inmemory_transactions

import (
	"context"
	"math/rand"
	"slices"
	"time"
	"user-wallet-service/internal/domain"
)

type R struct {
	state map[int][]domain.Transaction
}

func New() *R {
	return &R{
		state: make(map[int][]domain.Transaction, 0),
	}
}

func (r *R) Get(ctx context.Context, transactionID int) (*domain.Transaction, error) {
	for _, u := range r.state {
		for i := range u {
			if u[i].ID == transactionID {
				t := &domain.Transaction{
					ID:          u[i].ID,
					UserID:      u[i].UserID,
					Status:      u[i].Status,
					Amount:      u[i].Amount,
					Description: u[i].Description,
				}

				// rest a little bit in order to create race conditions
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				return t, nil
			}
		}
	}

	return nil, domain.NewErrTxNotFound(transactionID)
}

func (r *R) Create(ctx context.Context, tx domain.Transaction) (transactionID int, e error) {
	id := rand.Intn(9999)
	_, ok := r.state[tx.UserID]
	if !ok {
		r.state[tx.UserID] = make([]domain.Transaction, 1)
		r.state[tx.UserID][0] = domain.Transaction{
			ID:          id,
			UserID:      tx.UserID,
			Status:      tx.Status,
			Amount:      tx.Amount,
			Description: tx.Description,
		}

		return id, nil
	}

	r.state[tx.UserID] = append(r.state[tx.UserID], domain.Transaction{
		ID:     id,
		UserID: tx.UserID,
		Status: tx.Status,
		Amount: tx.Amount,
	})

	return id, nil
}

func (r *R) Total(ctx context.Context, userID int) (int, error) {
	var total int
	for _, tx := range r.state[userID] {
		if tx.Status == domain.TransactionStatusComplete || tx.Status == domain.TransactionStatusHold {
			total += tx.Amount
		}
	}

	return total, nil
}

func (r *R) Change(ctx context.Context, transactionID int, status domain.TransactionStatus) error {
	for _, u := range r.state {
		for i := range u {
			if u[i].ID == transactionID && slices.Contains(domain.TransactionStateMachine[u[i].Status], status) {
				u[i].Status = status
				return nil
			}
		}
	}

	return domain.NewErrTxNotFound(transactionID)
}
