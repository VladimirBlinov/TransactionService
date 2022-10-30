package sqlstore

import (
	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type TransactionRepo struct {
	store *Store
}

func (trr *TransactionRepo) Create(trm *model.Transaction) error {
	if err := trm.Validate(); err != nil {
		return err
	}

	tx, err := trr.store.db.Begin()
	if err != nil {
		return err
	}

	err = trr.store.db.QueryRow(
		`INSERT INTO public.transaction(
			amount, date_time)
			VALUES ($1, $2) RETURNING id`,
		trm.Amount,
		trm.DateTime,
	).Scan(&trm.ID)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = trr.store.db.QueryRow(
		`INSERT INTO public.user_transaction(
			user_id, transaction_id)
			VALUES ($1, $2)`,
		trm.UserID,
		trm.ID,
	).Err()

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (trr *TransactionRepo) GetBalance(*store.UserID) (*store.UserBalance, error) {
	return nil, nil
}
