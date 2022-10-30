package service

import "github.com/VladimirBlinov/TransactionService/Backend/internal/store"

type Service struct {
	TransactionService *TransactionService
	AuthService        *AuthService
	BalanceService     *BalanceService
}

func NewService(store store.Store) *Service {
	TransactionService := NewTransactionService(store)
	AuthService := NewAuthService(store)
	BalanceService := NewBalanceService(store)
	return &Service{
		TransactionService: TransactionService,
		AuthService:        AuthService,
		BalanceService:     BalanceService,
	}
}
