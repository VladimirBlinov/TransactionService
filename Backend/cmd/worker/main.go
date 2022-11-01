package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/app/workerserver"
	"github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	logger := logrus.New()

	config := workerserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logger.Fatalf("Failed load config", err.Error())
	}

	///////

	wrkserver, err := workerserver.NewWorkerServer(config)
	if err != nil {
		logger.Fatalf("Failed create WorkerServer", err.Error())
	}

	go func() {
		if err = wrkserver.Start(); err != nil {
			logger.Fatalf("error on worker server start", err.Error())
		}
	}()

	logger.Println("Started worker server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Println("Shutting down worker server...")
	if err = wrkserver.ShutDown(context.Background()); err != nil {
		logger.Fatalf("error on server shutdown: %s", err.Error())
	}
	os.Exit(0)

	////////

	// db, err := newDB(config.DataBaseURL)
	// failOnError(err, "Failed connect to DB")

	// defer func(db *sql.DB) {
	// 	if err = db.Close(); err != nil {
	// 		failOnError(err, "Error db close")
	// 	}
	// }(db)

	// store := sqlstore.New(db)
	// services := service.NewService(store)

	// rmq, err := rabbit.NewRabbitMQ(config.RabbitURL)
	// failOnError(err, "RabbitMQ init error")
	// defer rmq.Close()

	// activeUsersQueue, err := rmq.Channel.QueueDeclare(
	// 	"active_users", // name
	// 	true,           // durable
	// 	false,          // delete when unused
	// 	false,          // exclusive
	// 	false,          // no-wait
	// 	nil,            // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// activeUsersMsgs, err := rmq.Channel.Consume(
	// 	activeUsersQueue.Name, // queue
	// 	"",                    // consumer
	// 	true,                  // auto-ack
	// 	false,                 // exclusive
	// 	false,                 // no-local
	// 	false,                 // no-wait
	// 	nil,                   // args
	// )
	// failOnError(err, "Failed to register a consumer")

	// var activeUsers chan struct{}

	// go func() {
	// 	for activeUsersMsg := range activeUsersMsgs {
	// 		log.Printf("active_users msg[x] %s", activeUsersMsg.Body)

	// 		activeUsersTask := &service.InputTransaction{}
	// 		if err := json.Unmarshal(activeUsersMsg.Body, activeUsersTask); err != nil {
	// 			failOnError(err, "Failed to read active_users msg")
	// 		}

	// 		userQueue, err := rmq.Channel.QueueDeclare(
	// 			fmt.Sprintf("user.%d", activeUsersTask.UserID), // name
	// 			true,  // durable
	// 			false, // delete when unused
	// 			false, // exclusive
	// 			false, // no-wait
	// 			nil,   // arguments
	// 		)
	// 		failOnError(err, "Failed to declare a queue")

	// 		userMsgs, err := rmq.Channel.Consume(
	// 			userQueue.Name, // queue
	// 			"",             // consumer
	// 			true,           // auto-ack
	// 			false,          // exclusive
	// 			false,          // no-local
	// 			false,          // no-wait
	// 			nil,            // args
	// 		)
	// 		failOnError(err, "Failed to register a consumer")

	// 		go func(userQueue string) {
	// 			for userMsg := range userMsgs {
	// 				log.Printf("%s [x] %s", userQueue, userMsg.Body)

	// 				tw := worker.NewTransactionWorker(services)
	// 				resp := tw.Run(userMsg.Body)
	// 				log.Printf("Response %s [x] %s", userQueue, resp)
	// 			}
	// 		}(userQueue.Name)

	// 		log.Printf(" [*] Waiting for user transaction. To exit press CTRL+C")
	// 	}
	// }()

	// log.Printf("active_users [*] Waiting for users. To exit press CTRL+C")
	// <-activeUsers
}

// func newDB(databaseURL string) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", databaseURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = db.Ping(); err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
