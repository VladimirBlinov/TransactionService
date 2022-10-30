package store

type UserID int32

type TransactionRepo interface {
	Create() error
	GetBalance(*UserID) (float32, error)
}
