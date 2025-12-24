package events

import "time"

const TopicRefundUpdate = "refund.update"
const TopicAddAttempts = "add.attempts"

type RefundUpdateV1 struct {
	OperationID string    `json:"operation_id"`
	TestTakerID uint64    `json:"test_taker_id"`
	Variant     int32     `json:"variant"`
	Lang        string    `json:"lang"`
	Refund      int32     `json:"refund"`
	OccurredAt  time.Time `json:"occurred_at"`
}
type AddAttemptsV1 struct {
	OperationID string    `json:"operation_id"`
	TestTakerID uint64    `json:"test_taker_id"`
	Variant     int32     `json:"variant"`
	Lang        string    `json:"lang"`
	Attempts    int32     `json:"attempts"`
	OccurredAt  time.Time `json:"occurred_at"`
}
