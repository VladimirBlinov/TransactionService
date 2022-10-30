package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
	}
}

func TestTransaction(t *testing.T) *Transaction {
	return &Transaction{
		UserID:   1,
		Amount:   100,
		DateTime: time.Now(),
	}
}

func TestBalance(t *testing.T) *Balance {
	return &Balance{
		UserID:        1,
		Balance:       100,
		AuditDateTime: time.Now(),
	}
}
