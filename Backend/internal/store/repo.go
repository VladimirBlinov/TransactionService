package store

import "github.com/VladimirBlinov/TransactionService/Backend/internal/model"

type TransactionRepo interface {
	Create(*model.Transaction) error
}

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type BalanceRepo interface {
	Create(*model.Balance) error
	GetBalanceByUserID(int) (*model.Balance, error)
	UpdateBalance(*model.Balance) error
}
