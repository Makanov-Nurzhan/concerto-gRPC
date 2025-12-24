package domain

import "errors"

var (
	ErrInvalidRefund             = errors.New("invalid refund")
	ErrInvalidAttemptCount       = errors.New("invalid attempt count")
	ErrInvalidTotalAttempts      = errors.New("invalid total attempts")
	ErrFirstDayAttempts          = errors.New("first day attempts")
	ErrHasActiveSession          = errors.New("user has active session")
	ErrFailedToUpdate            = errors.New("failed to update")
	ErrOperationAlreadyProcessed = errors.New("operation already processed")
	ErrOperationInProgress       = errors.New("operation is in progress")
	ErrInvalidAttemptToAdd       = errors.New("invalid attempt to add")
)
