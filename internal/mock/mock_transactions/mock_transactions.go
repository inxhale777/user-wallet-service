package mock_transactions

import (
	"context"
	"github.com/pkg/errors"
	"math/rand"
	"user-balance-service/internal/domain"
)

type R struct {
	state map[int][]domain.Transaction
}

func New() *R {
	return &R{}
}

func (r *R) Get(ctx context.Context, transactionID int) (*domain.Transaction, error) {
	for _, u := range r.state {
		for i := range u {
			if u[i].ID == transactionID {
				return &domain.Transaction{
					ID:        u[i].ID,
					UserID:    u[i].UserID,
					ServiceID: u[i].ServiceID,
					OrderID:   u[i].OrderID,
					Status:    u[i].Status,
					Amount:    u[i].Amount,
				}, nil
			}
		}
	}

	return nil, errors.New("transaction not found")
}

func (r *R) Create(ctx context.Context, userID int, amount int, status domain.TransactionStatus) (transactionID int, e error) {
	id := rand.Intn(9999)
	_, ok := r.state[userID]
	if !ok {
		r.state[userID] = make([]domain.Transaction, 1)
		r.state[userID][0] = domain.Transaction{
			ID:     id,
			UserID: userID,
			Status: status,
			Amount: amount,
		}

		return id, nil
	}

	r.state[userID] = append(r.state[userID], domain.Transaction{
		ID:     id,
		UserID: userID,
		Status: status,
		Amount: amount,
	})

	return id, nil
}

func (r *R) Total(ctx context.Context, userID int) (int, error) {
	var total int
	for _, tx := range r.state[userID] {
		total += tx.Amount
	}

	return total, nil
}

func (*R) Change(ctx context.Context, transactionID int, status domain.TransactionStatus) error {
	//TODO implement me
	panic("implement me")
}
