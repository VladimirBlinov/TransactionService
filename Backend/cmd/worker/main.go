package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/app/worker"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/service"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/store/sqlstore"

	amqp "github.com/rabbitmq/amqp091-go"
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
	if err != nil {
		failOnError(err, "Failed load config")
	}

	db, err := newDB(config.DataBaseURL)
	if err != nil {
		failOnError(err, "Failed connect to DB")
	}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			failOnError(err, "Error db close")
		}
	}(db)

	store := sqlstore.New(db)
	services := service.NewService(store)
	tw := worker.NewTransactionWorker(services)

	rmq, err := rabbit.NewRabbitMQ()
	if err != nil {
		return nil, fmt.Errorf("internal.NewRabbitMQ %w", err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
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

			tw.Run(d.Body)
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
