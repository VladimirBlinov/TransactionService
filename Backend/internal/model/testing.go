package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "ex@test.org",
		Password: "password",
		Balance:  0,
	}
}

func TestTransaction(t *testing.T) *Transaction {
	return &Transaction{
		UserID:   1,
		Amount:   100,
		DateTime: time.Now(),
	}
}
