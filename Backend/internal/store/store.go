package store

type Store interface {
	Transaction() TransactionRepo
	User() UserRepo
}
