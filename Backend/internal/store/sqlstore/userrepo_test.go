package sqlstore_test

import (
	"testing"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
	assert.NotEqual(t, 0, u.BalanceID)
}

func TestFindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)

	email := "user@example.org"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestFindById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	s.User().Create(u1)

	u2, err := s.User().FindById(u1.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func Test_GetBalance(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "balance", "balance_audit", "user_balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	u.Balance = 200
	s.User().Create(u)

	b, err := s.User().GetBalance(u.ID)

	assert.NoError(t, err)
	assert.Equal(t, u.Balance, b)
}
