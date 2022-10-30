package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Transaction struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"date_time"`
	IsValid  bool
}

func (tr *Transaction) Validate() error {
	return validation.ValidateStruct(
		tr,
		validation.Field(&tr.UserID, validation.Required),
		validation.Field(&tr.Amount, validation.Required),
		validation.Field(&tr.DateTime, validation.Required),
	)
}

func (tr *Transaction) CheckIsValid(balance float64) error {
	if balance+tr.Amount >= 0 {
		tr.IsValid = true
		return nil
	}

	tr.IsValid = false

	return ErrTransactionOutOfBalance
}
