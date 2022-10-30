package model_test

import (
	"testing"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_BalanceValidate(t *testing.T) {
	testCases := []struct {
		name    string
		b       func() *model.Balance
		isValid bool
	}{
		{
			name: "valid",
			b: func() *model.Balance {
				return model.TestBalance(t)
			},
			isValid: true,
		},
		{
			name: "balance equal zero",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.Balance = 0
				return b
			},
			isValid: true,
		},
		{
			name: "empty user_id",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.UserID = 0
				return b
			},
			isValid: false,
		},
		{
			name: "empty date_time",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.AuditDateTime = time.Time{}
				return b
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.b().Validate())
			} else {
				assert.Error(t, tc.b().Validate())
			}
		})
	}
}
