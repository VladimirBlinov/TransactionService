package apiserver_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5435 user=admin password=qaz dbname=TransactionTest sslmode=disable"
	}

	os.Exit(m.Run())
}
