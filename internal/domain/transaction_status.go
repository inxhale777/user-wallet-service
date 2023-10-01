package domain

type TransactionStatus string

var (
	TransactionStatusHold      TransactionStatus = "hold"
	TransactionStatusComplete  TransactionStatus = "complete"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

var TransactionStateMachine = map[TransactionStatus][]TransactionStatus{
	// from hold to complete or cancelled
	TransactionStatusHold: {TransactionStatusComplete, TransactionStatusCancelled},
	// from cancelled to nothing
	TransactionStatusComplete: {},
	// from complete to nothing
	TransactionStatusCancelled: {},
}
