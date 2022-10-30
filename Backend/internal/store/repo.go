package store

import "github.com/VladimirBlinov/TransactionService/Backend/internal/model"

type UserID int32
type UserBalance float32

type TransactionRepo interface {
	Create(*model.Transaction) error
	GetBalance(*UserID) (*UserBalance, error)
}

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
