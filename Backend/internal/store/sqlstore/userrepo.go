package sqlstore

import (
	"database/sql"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
)

type UserRepo struct {
	store *Store
}

func (r *UserRepo) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	err := u.EncryptPasswordBeforeCreate()
	if err != nil {
		return err
	}

	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}

	err = r.store.db.QueryRow(
		"INSERT INTO public.users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = r.store.db.QueryRow(
		"INSERT INTO public.balance (active) VALUES ($1) RETURNING id",
		true,
	).Scan(&u.BalanceID)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = r.store.db.QueryRow(
		"INSERT INTO public.user_balance (user_id, balance_id) VALUES ($1, $2)",
		u.ID,
		u.BalanceID,
	).Err()

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	err = r.store.db.QueryRow(
		"INSERT INTO public.balance_audit (balance_id, balance, last_audit_time) VALUES ($1, $2, $3)",
		u.BalanceID,
		u.Balance,
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

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM public.users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM public.users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) GetBalance(userID int) (float64, error) {
	balance := 0.0
	if err := r.store.db.QueryRow(
		`SELECT balance FROM public.balance_audit as ba
		inner join public.balance as b on b.id = ba.balance_id
		inner join public.user_balance as ub on b.id = ub.balance_id
		
		where b.active = true and ub.user_id = $1
		order by ba.last_audit_time desc
		limit 1`,
		userID,
	).Scan(
		&balance,
	); err != nil {
		if err == sql.ErrNoRows {
			return balance, store.ErrRecordNotFound
		}
		return balance, err
	}

	return balance, nil
}
