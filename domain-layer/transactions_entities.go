package models

type TransactionsEntity struct {
	ID      int
	Entity  string
	Kind    string
	Balance float32
}

type TransactionsEntities []TransactionsEntity

var (
	RefToTransactionsEntities *TransactionsEntities
)
