package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID              int
	TransactionDate time.Time
	Transaction     string
	Categories      Categories
	Kind            string
	Amount          float32
}

type Transactions []Transaction

type TransactionEntity string
type TransactionKind string

const (
	TestEntity                 TransactionEntity = "test"
	Test2Entity                TransactionEntity = "test2"
	DebitBankAccountKind       TransactionKind   = "debit_bank_account"
	DebitCreditBankAccountKind TransactionKind   = "debit_credit_bank_account"
	DebitKindTransaction       string            = "debit"
	CreditKindTransaction      string            = "credit"
)

var (
	TransactionEntities = []TransactionEntity{
		TestEntity, Test2Entity,
	}

	TransactionKinds = []TransactionKind{
		DebitBankAccountKind, DebitCreditBankAccountKind,
	}

	RefsToTransactions = &map[string]*Transactions{
		EntityKindKey(TestEntity, DebitBankAccountKind):        nil,
		EntityKindKey(Test2Entity, DebitCreditBankAccountKind): nil,
	}
)

func EntityKindKey(entity TransactionEntity, kind TransactionKind) string {
	return fmt.Sprintf("%s_#_%s", entity, kind)
}
