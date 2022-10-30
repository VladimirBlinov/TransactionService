package service

import "github.com/VladimirBlinov/TransactionService/Backend/internal/store"

type TransactionService struct {
	store store.Store
}

func NewTransactionService(store store.Store) *TransactionService {
	return &TransactionService{
		store: store,
	}
}
