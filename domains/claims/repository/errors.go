package repository

import "errors"

var (
	ErrClaimNumberCollision = errors.New("claim number collision")

	ErrClaimFKViolation = errors.New("claim foreign key violation")
)

type ClaimFKViolationError struct {
	Constraint string
	Detail     string
}

func (e *ClaimFKViolationError) Error() string {
	if e.Detail != "" {
		return "claim foreign key violation: " + e.Constraint + ": " + e.Detail
	}
	return "claim foreign key violation: " + e.Constraint
}

func (e *ClaimFKViolationError) Unwrap() error {
	return ErrClaimFKViolation
}
