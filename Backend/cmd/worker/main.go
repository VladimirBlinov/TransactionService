package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/app/worker"
	rabbit "github.com/VladimirBlinov/TransactionService/Backend/internal/rabbitmq"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"
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
	config := worker.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	failOnError(err, "Failed load config")

	db, err := newDB(config.DataBaseURL)
	failOnError(err, "Failed connect to DB")

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			failOnError(err, "Error db close")
		}
	}(db)

	store := sqlstore.New(db)
	services := service.NewService(store)
	tw := worker.NewTransactionWorker(services)

	rmq, err := rabbit.NewRabbitMQ()
	failOnError(err, "RabbitMQ init error")
	defer rmq.Close()

	q, err := rmq.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = rmq.Channel.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"users", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := rmq.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)

			resp := tw.Run(d.Body)
			log.Printf("Response [x] %s", resp)
		}
	}()

	log.Printf(" [*] Waiting for transaction. To exit press CTRL+C")
	<-forever
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
