package worker

import (
	"encoding/json"
	"time"

	"github.com/VladimirBlinov/TransactionService/Backend/internal/model"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
)

type TransactionWorker struct {
	service *service.Service
}

func NewTransactionWorker(service *service.Service) *TransactionWorker {
	return &TransactionWorker{
		service: service,
	}
}

func (tw *TransactionWorker) Run(task []byte) []byte {
	tr := &model.Transaction{}

	if err := json.Unmarshal(task, tr); err != nil {
		return tw.error(err)
	}

	tr.DateTime = time.Now()

	err := tw.service.TransactionService.CreateTransaction(tr)
	if err != nil {
		return tw.error(err)
	}

	u, err := tw.service.AuthService.Authenticate(tr.UserID)
	if err != nil {
		return tw.error(err)
	}

	b, err := tw.service.BalanceService.ApplyTransaction(u, tr)
	if err != nil {
		return tw.error(err)
	}

	return tw.respond(b)
}

func (tw *TransactionWorker) error(err error) []byte {
	return tw.respond(map[string]string{"error": err.Error()})
}

func (tw *TransactionWorker) respond(data interface{}) []byte {
	response, err := json.Marshal(data)
	if err != nil {
		tw.error(err)
	}
	return response
}
