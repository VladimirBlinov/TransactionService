package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(connUrl string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(connUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// err = ch.ExchangeDeclare(
	// 	"",    // name
	// 	"",    // type
	// 	true,  // durable
	// 	false, // auto-deleted
	// 	false, // internal
	// 	false, // no-wait
	// 	nil,   // arguments
	// )
	// failOnError(err, "Failed to declare an exchange")

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close ...
func (r *RabbitMQ) Close() {
	r.Connection.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
