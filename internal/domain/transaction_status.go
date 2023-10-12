package domain

type TransactionStatus string

var (
	TransactionStatusHold      TransactionStatus = "hold"
	TransactionStatusComplete  TransactionStatus = "complete"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// TransactionStateMachine "wanted status": "need to have status".
var TransactionStateMachine = map[TransactionStatus]TransactionStatus{
	// unable to change status to HOLD
	TransactionStatusHold: "",
	// able to change status to COMPLETE only from HOLD
	TransactionStatusComplete: TransactionStatusHold,
	// able to change status to CANCELLED only from HOLD
	TransactionStatusCancelled: TransactionStatusHold,
}
