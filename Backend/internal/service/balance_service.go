package service

import (
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type BalanceService struct {
	store store.Store
}

func NewBalanceService(store store.Store) *BalanceService {
	return &BalanceService{
		store: store,
	}
}

func (bs *BalanceService) ApplyTransaction(u *model.User, tr *model.Transaction) (*model.Balance, error) {
	b, err := bs.store.Balance().GetBalanceByUserID(u.ID)
	if err != nil {
		return b, err
	}

	if err = tr.CheckIsValid(b.Balance); err != nil {
		return b, err
	}

	b.ChangeBalance(tr.Amount, tr.DateTime)

	if err = bs.store.Balance().UpdateBalance(b); err != nil {
		return b, err
	}

	return b, nil
}

func (bs *BalanceService) CreateBalance(u *model.User) (*model.Balance, error) {
	b := &model.Balance{}
	b.UserID = u.ID
	b.Balance = 0
	b.AuditTime = time.Now()

	if err := bs.store.Balance().Create(b); err != nil {
		return nil, err
	}

	return b, nil
}
