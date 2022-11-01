package apiserver_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
	rabbitURL   string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5435 user=admin password=qaz dbname=TransactionTest sslmode=disable"
	}

	rabbitURL = os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	os.Exit(m.Run())
}
