package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Balance struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Balance       float32   `json:"balance"`
	AuditDateTime time.Time `json:"audit_date_time"`
}

func (b *Balance) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.UserID, validation.Required),
	)
}
