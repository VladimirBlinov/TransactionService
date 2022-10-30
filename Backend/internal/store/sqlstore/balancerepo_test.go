package sqlstore_test

import (
	"testing"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_BalanceCreate(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	b := model.TestBalance(t)
	b.UserID = u.ID

	err := s.Balance().Create(b)

	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func Test_GetBalanceByUserID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	b := model.TestBalance(t)
	b.UserID = u.ID
	s.Balance().Create(b)

	ub, err := s.Balance().GetBalanceByUserID(u.ID)

	assert.NoError(t, err)
	assert.Equal(t, b.Balance, ub.Balance)
}

func Test_UpdateBalance(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	b := model.TestBalance(t)
	b.UserID = u.ID
	s.Balance().Create(b)

	initBalance := b.Balance

	tr := model.TestTransaction(t)

	b.ChangeBalance(tr.Amount, tr.DateTime)

	assert.NoError(t, s.Balance().UpdateBalance(b))

	bu, _ := s.Balance().GetBalanceByUserID(u.ID)
	assert.Equal(t, initBalance+tr.Amount, bu.Balance)
}
