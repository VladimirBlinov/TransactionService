package workerserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	rabbit "github.com/VladimirBlinov/TransactionService/Backend/internal/rabbitmq"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/worker"
	"github.com/sirupsen/logrus"
)

type WorkerServer struct {
	db       *sql.DB
	store    store.Store
	services service.Service
	rmq      rabbit.RabbitMQ
	logger   logrus.Logger
}

func NewWorkerServer(config *Config) (*WorkerServer, error) {
	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return nil, errors.New("Failed connect to DB")
	}

	// defer func(db *sql.DB) {
	// 	if err = db.Close(); err != nil {
	// 		logrus.Errorf("error db close: %s", err.Error())
	// 	}
	// }(db)

	store := sqlstore.New(db)
	services := service.NewService(store)

	rmq, err := rabbit.NewRabbitMQ(config.RabbitURL)
	if err != nil {
		return nil, errors.New("RabbitMQ init error")
	}

	//defer rmq.Close()

	return (&WorkerServer{
		db:       db,
		store:    store,
		services: *services,
		rmq:      *rmq,
		logger:   *logrus.New(),
	}), nil
}

func (ws *WorkerServer) Start() error {
	activeUsersQueue, err := ws.rmq.Channel.QueueDeclare(
		"active_users", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		ws.logger.Fatalf("ActiveUsersQueue error", err.Error())
		return err
	}

	activeUsersMsgs, err := ws.rmq.Channel.Consume(
		activeUsersQueue.Name, // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		ws.logger.Fatalf("activeUsersMsgs error", err.Error())
		return err
	}

	var activeUsers chan struct{}

	go func() {
		for activeUsersMsg := range activeUsersMsgs {
			ws.logger.Printf("active_users msg[x] %s", activeUsersMsg.Body)

			activeUsersTask := &service.InputTransaction{}
			if err := json.Unmarshal(activeUsersMsg.Body, activeUsersTask); err != nil {
				ws.logger.Fatalf("Failed to read active_users msg", err.Error())
			}

			userQueue, err := ws.rmq.Channel.QueueDeclare(
				fmt.Sprintf("user.%d", activeUsersTask.UserID), // name
				true,  // durable
				false, // delete when unused
				false, // exclusive
				false, // no-wait
				nil,   // arguments
			)
			if err != nil {
				ws.logger.Fatalf("Failed to declare a queue", err.Error())
			}

			userMsgs, err := ws.rmq.Channel.Consume(
				userQueue.Name, // queue
				"",             // consumer
				true,           // auto-ack
				false,          // exclusive
				false,          // no-local
				false,          // no-wait
				nil,            // args
			)
			if err != nil {
				ws.logger.Fatalf("Failed to register a consumer", err.Error())
			}

			go func(userQueue string) {
				for userMsg := range userMsgs {
					ws.logger.Printf("%s [x] %s", userQueue, userMsg.Body)

					tw := worker.NewTransactionWorker(&ws.services)
					resp := tw.Run(userMsg.Body)
					ws.logger.Printf("Response %s [x] %s", userQueue, resp)
				}
			}(userQueue.Name)

			ws.logger.Printf(" [*] Waiting for user transaction. To exit press CTRL+C")
		}
	}()

	ws.logger.Printf("active_users [*] Waiting for users. To exit press CTRL+C")
	<-activeUsers

	return nil
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (ws *WorkerServer) ShutDown(ctx context.Context) error {
	ws.rmq.Close()
	if err := ws.db.Close(); err != nil {
		logrus.Errorf("error db close: %s", err.Error())
	}
	return nil
}
