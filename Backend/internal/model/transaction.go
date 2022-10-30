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
}

func (tr *Transaction) Validate() error {
	return validation.ValidateStruct(
		tr,
		validation.Field(&tr.UserID, validation.Required),
		validation.Field(&tr.Amount, validation.Required),
		validation.Field(&tr.DateTime, validation.Required),
	)
}
