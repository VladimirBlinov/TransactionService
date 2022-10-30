package store

import "github.com/VladimirBlinov/TransactionService/Backend/internal/model"

type TransactionRepo interface {
	Create(*model.Transaction) error
	//GetBalance(int) (float32, error)
}

type UserRepo interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
