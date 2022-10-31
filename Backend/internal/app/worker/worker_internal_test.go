package worker_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/app/worker"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func Test_WorkerRun(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "transaction", "user_transaction", "user_balance", "balance_audit", "balance")

	store := sqlstore.New(db)
	services := service.NewService(store)
	tw := worker.NewTransactionWorker(services)

	u := model.TestUser(t)
	store.User().Create(u)

	b := model.TestBalance(t)
	b.UserID = u.ID
	store.Balance().Create(b)

	testCases := []struct {
		name    string
		message map[string]interface{}
		success bool
	}{
		{
			name: "valid",
			message: map[string]interface{}{
				"user_id":   u.ID,
				"amount":    200.0,
				"date_time": time.Now(),
			},
			success: true,
		},
		{
			name: "invalid",
			message: map[string]interface{}{
				"user_id":   u.ID,
				"amount":    -1200.0,
				"date_time": time.Now(),
			},
			success: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.message)
			if err != nil {
				return
			}
			resp := tw.Run(b)
			resp_str := string(resp[:])
			if tc.success {
				assert.NotContains(t, resp_str, "error")
			} else {
				assert.Contains(t, resp_str, "error")
			}
		})
	}
}
