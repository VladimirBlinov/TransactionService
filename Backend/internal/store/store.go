package store

type Store interface {
	Transaction() TransactionRepo
}
