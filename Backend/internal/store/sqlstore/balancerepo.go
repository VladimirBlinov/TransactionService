package sqlstore

import (
	"database/sql"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type BalanceRepo struct {
	store *Store
}

func (br *BalanceRepo) Create(b *model.Balance) error {
	if err := b.Validate(); err != nil {
		return err
	}

	tx, err := br.store.db.Begin()
	if err != nil {
		return err
	}

	err = br.store.db.QueryRow(
		"INSERT INTO public.balance (active) VALUES ($1) RETURNING id",
		true,
	).Scan(&b.ID)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = br.store.db.QueryRow(
		"INSERT INTO public.user_balance (user_id, balance_id) VALUES ($1, $2)",
		b.UserID,
		b.ID,
	).Err()

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = br.store.db.QueryRow(
		"INSERT INTO public.balance_audit (balance_id, balance, last_audit_time) VALUES ($1, $2, $3)",
		b.ID,
		b.Balance,
		time.Now(),
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

func (br *BalanceRepo) GetBalanceByUserID(userID int) (*model.Balance, error) {
	b := &model.Balance{}
	if err := br.store.db.QueryRow(
		`SELECT b.id, ba.balance, ub.user_id, ba.last_audit_time FROM public.balance_audit as ba
		inner join public.balance as b on b.id = ba.balance_id
		inner join public.user_balance as ub on b.id = ub.balance_id
		
		where b.active = true and ub.user_id = $1
		order by ba.last_audit_time desc
		limit 1`,
		userID,
	).Scan(
		&b.ID,
		&b.Balance,
		&b.UserID,
		&b.AuditTime,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrBalanceRecordNotFound
		}
		return nil, err
	}

	return b, nil
}

func (br *BalanceRepo) UpdateBalance(b *model.Balance) error {
	if err := b.Validate(); err != nil {
		return err
	}

	err := br.store.db.QueryRow(
		"INSERT INTO public.balance_audit (balance_id, balance, last_audit_time) VALUES ($1, $2, $3)",
		b.ID,
		b.Balance,
		b.AuditTime,
	).Err()

	if err != nil {
		return err
	}

	return nil
}
