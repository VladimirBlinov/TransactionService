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
	u.Balance = 0

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
