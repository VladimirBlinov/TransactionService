package sqlstore

import (
	"database/sql"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db              *sql.DB
	userRepo        *UserRepo
	transactionRepo *TransactionRepo
	balanceRepo     *BalanceRepo
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepo {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = &UserRepo{
		store: s,
	}
	return s.userRepo
}

func (s *Store) Balance() store.BalanceRepo {
	if s.balanceRepo != nil {
		return s.balanceRepo
	}

	s.balanceRepo = &BalanceRepo{
		store: s,
	}
	return s.balanceRepo
}

func (s *Store) TransactionRepo() store.TransactionRepo {
	if s.transactionRepo != nil {
		return s.transactionRepo
	}

	s.transactionRepo = &TransactionRepo{
		store: s,
	}
	return s.transactionRepo
}
