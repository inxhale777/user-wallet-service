package domain

import "errors"

var (
	ErrNoLockerProvided = errors.New("you must provide locker impl. to use this method")
	ErrInvalidAmount    = errors.New("amount must be greater or equal 1")
)

type GenericError struct {
	Message string
}

func (e *GenericError) Error() string {
	return e.Message
}

type ErrTxNotFound struct {
	GenericError
	TransactionID int
}

func NewErrTxNotFound(transactionID int) *ErrTxNotFound {
	return &ErrTxNotFound{
		GenericError: GenericError{
			Message: "transaction not found",
		},
		TransactionID: transactionID,
	}
}

type ErrInsufficientMoney struct {
	GenericError
	UserID int
	Want   int
	Have   int
}

func NewErrInsufficientMoney(userID, want, have int) *ErrInsufficientMoney {
	return &ErrInsufficientMoney{
		GenericError: GenericError{
			Message: "insufficient money",
		},
		UserID: userID,
		Want:   want,
		Have:   have,
	}
}

type ErrInvalidTxStatus struct {
	GenericError
	TransactionID int
	Want          TransactionStatus
	Have          TransactionStatus
}

func NewErrInvalidTxStatus(transactionID int, want, have TransactionStatus) *ErrInvalidTxStatus {
	return &ErrInvalidTxStatus{
		GenericError: GenericError{
			Message: "invalid tx status",
		},
		TransactionID: transactionID,
		Want:          want,
		Have:          have,
	}
}
