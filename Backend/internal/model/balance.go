package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Balance struct {
	ID        int       `json:"id"`
	Balance   float64   `json:"balance"`
	UserID    int       `json:"user_id"`
	AuditTime time.Time `json:"audit_time"`
}

func (b *Balance) ChangeBalance(amount float64, auditTime time.Time) error {
	b.Balance += amount
	b.AuditTime = auditTime
	return nil
}

func (b *Balance) ValidateBalanceChange(amount float64) bool {
	if b.Balance+amount >= 0 {
		return true
	}
	return false
}

func (b *Balance) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.Balance, validation.Min(float64(0))),
		validation.Field(&b.UserID, validation.Required),
	)
}
