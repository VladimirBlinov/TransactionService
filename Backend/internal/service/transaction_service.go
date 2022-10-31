package service

import (
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type InputTransaction struct {
	UserID   int       `json:"user_id"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"date_time"`
}

type TransactionService struct {
	store store.Store
}

func NewTransactionService(store store.Store) *TransactionService {
	return &TransactionService{
		store: store,
	}
}

func (trs *TransactionService) CreateTransaction(tr *model.Transaction) error {
	if err := trs.store.Transaction().Create(tr); err != nil {
		return err
	}

	return nil
}
