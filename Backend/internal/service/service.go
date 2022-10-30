package service

import "github.com/VladimirBlinov/TransactionService/Backend/internal/store"

type Service struct {
	TransactionService *TransactionService
}

func NewService(store store.Store) *Service {
	TransactionService := NewTransactionService(store)
	return &Service{
		TransactionService: TransactionService,
	}
}
