package sqlstore_test

import (
	"testing"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepo_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	_ = s.User().Create(u)

	tr := model.TestTransaction(t)
	tr.UserID = u.ID
	assert.NoError(t, s.Transaction().Create(tr))
	assert.NotNil(t, tr)
}
