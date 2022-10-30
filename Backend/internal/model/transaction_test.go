package model_test

import (
	"testing"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionValidate(t *testing.T) {
	testCases := []struct {
		name    string
		tr      func() *model.Transaction
		isValid bool
	}{
		{
			name: "valid",
			tr: func() *model.Transaction {
				return model.TestTransaction(t)
			},
			isValid: true,
		},
		{
			name: "empty amount",
			tr: func() *model.Transaction {
				tr := model.TestTransaction(t)
				tr.Amount = 0
				return tr
			},
			isValid: false,
		},
		{
			name: "empty user_id",
			tr: func() *model.Transaction {
				tr := model.TestTransaction(t)
				tr.UserID = 0
				return tr
			},
			isValid: false,
		},
		{
			name: "empty date_time",
			tr: func() *model.Transaction {
				tr := model.TestTransaction(t)
				tr.DateTime = time.Time{}
				return tr
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.tr().Validate())
			} else {
				assert.Error(t, tc.tr().Validate())
			}
		})
	}
}

func Test_TransactionCheckValid(t *testing.T) {
	testCases := []struct {
		name    string
		tr      func() *model.Transaction
		balance float64
		isValid bool
	}{
		{
			name: "valid",
			tr: func() *model.Transaction {
				return model.TestTransaction(t)
			},
			balance: 100,
			isValid: true,
		},
		{
			name: "invalid",
			tr: func() *model.Transaction {
				tr := model.TestTransaction(t)
				tr.Amount = -200
				return tr
			},
			balance: 100,
			isValid: false,
		},
		{
			name: "valid balance to zero",
			tr: func() *model.Transaction {
				tr := model.TestTransaction(t)
				tr.Amount = -200.0
				return tr
			},
			balance: 200.0,
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.tr().CheckIsValid(tc.balance))
			} else {
				assert.Error(t, tc.tr().CheckIsValid(tc.balance))
			}
		})
	}
}
