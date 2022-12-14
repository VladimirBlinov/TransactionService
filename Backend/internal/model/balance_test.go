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
			name: "invalid",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.Balance = -100
				return b
			},
			isValid: false,
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

func Test_BalanceChange(t *testing.T) {
	testCases := []struct {
		name   string
		b      func() *model.Balance
		amount float64
		expect float64
	}{
		{
			name: "valid",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.Balance = 100
				return b
			},
			amount: 100,
			expect: 200,
		},
		{
			name: "invalid",
			b: func() *model.Balance {
				b := model.TestBalance(t)
				b.Balance = 100
				return b
			},
			amount: -200,
			expect: -100,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := tc.b()
			b.ChangeBalance(tc.amount, time.Now())
			assert.Equal(t, tc.expect, b.Balance)
		})
	}
}
