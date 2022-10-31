package store

import "errors"

var (
	ErrUserRecordNotFound    = errors.New("User record not found")
	ErrBalanceRecordNotFound = errors.New("Balance record not found")
)
