package repositories

import (
	"fmt"
	"strconv"
	"strings"
models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
"github.com/h-abranches-dev/daily-expenses-be/utils"
"time"
)
type TransactionDAO struct {
	ID              int
	TransactionDate time.Time
	Transaction     string
	Categories      CategoriesDAO
	Kind            string
	Amount          float32
}

type TransactionsDAO []TransactionDAO

type TransactionsRepo struct {
	*Repo
	Kind   string
	Entity string
}

const (
	bankAccountDebitPattern       = "%d;%s;%s;%s;%s"
	bankAccountDebitCreditPattern = "%d;%s;%s;%s;%s;%s"
)

var (
	dBFiles = map[string]string{
		models.EntityKindKey(models.TestEntity, models.DebitBankAccountKind):        "db/test_transactions.csv",
		models.EntityKindKey(models.Test2Entity, models.DebitCreditBankAccountKind): "db/test2_transactions.csv",
	}
)

func NewTransactionsRepo(kind, entity string) (*TransactionsRepo, error) {
	dbFile := dBFiles[models.EntityKindKey(models.TransactionEntity(entity), models.TransactionKind(kind))]
	r, err := NewRepo(dbFile)
	if err != nil {
		return nil, err
	}
	return &TransactionsRepo{
		Repo:   r,
		Kind:   kind,
		Entity: entity,
	}, nil
}

func (repo TransactionsRepo) ToRow(t TransactionDAO) (string, error) {
	if strings.Index(t.Transaction, repo.FileSeparator) != -1 {
		return "", fmt.Errorf("invalid 'Transaction' because includes the char %q => %q", repo.FileSeparator,
			t.Transaction)
	}

	var rowPattern, row string
	switch models.TransactionKind(repo.Kind) {
	case models.DebitBankAccountKind:
		rowPattern = bankAccountDebitPattern
		row = fmt.Sprintf(rowPattern, t.ID, t.TransactionDate.Format(utils.DateFormat),
			t.Transaction, strings.Join(t.categoriesLabels(), "#"), fmt.Sprintf("%.2f", t.Amount))
	case models.DebitCreditBankAccountKind:
		rowPattern = bankAccountDebitCreditPattern
		row = fmt.Sprintf(rowPattern, t.ID, t.TransactionDate.Format(utils.DateFormat),
			t.Transaction, strings.Join(t.categoriesLabels(), "#"), t.Kind, fmt.Sprintf("%.2f", t.Amount))
	default:
		return "", fmt.Errorf("invalid repository kind")
	}

	return row, nil
}

func (repo TransactionsRepo) rowToTransaction(categoriesRepo CategoriesRepo, row string) (TransactionDAO, error) {
	emptyTransaction := TransactionDAO{}
	columns := strings.Split(row, ";")
	tID, err := strconv.Atoi(columns[0])
	if err != nil {
		return emptyTransaction, err
	}
	tDate, err := time.Parse(utils.DateFormat, strings.Trim(columns[1], " "))
	if err != nil {
		return emptyTransaction, err
	}

	kindColumnIdx := -1
	amountColumnIdx := -1
	var tKind string
	switch models.TransactionKind(repo.Kind) {
	case models.DebitBankAccountKind:
		amountColumnIdx = 4
	case models.DebitCreditBankAccountKind:
		kindColumnIdx = 4
		tKind = columns[kindColumnIdx]
		amountColumnIdx = 5
	default:
		return emptyTransaction, fmt.Errorf("invalid repository kind")
	}

	tAmount, err := strconv.ParseFloat(strings.Trim(columns[amountColumnIdx], " "), 32)
	if err != nil {
		return emptyTransaction, err
	}

	csDAO := CategoriesDAO{}
	if columns[3] != "" {
		csDAO, err = categoriesRepo.CategoriesDAO(strings.Split(columns[3], "#"))
		if err != nil {
			return emptyTransaction, err
		}
	}
	return TransactionDAO{
		ID:              tID,
		TransactionDate: tDate,
		Transaction:     columns[2],
		Categories:      csDAO,
		Kind:            tKind,
		Amount:          float32(tAmount),
	}, nil
}

func (repo TransactionsRepo) GetAllTransactions() (TransactionsDAO, error) {
	transactions := TransactionsDAO{}
	rows := repo.FileWrapper.Lines
	csRepo, err := GetCategoriesRepo()
	if err != nil {
		return transactions, err
	}

	for i := 1; i < len(*rows)-1; i++ {
		var transaction TransactionDAO
		transaction, err = repo.rowToTransaction(*csRepo, (*rows)[i])
		if err != nil {
			return TransactionsDAO{}, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (repo TransactionsRepo) AddTransaction(t TransactionDAO) error {
	line, err := repo.ToRow(t)
	if err != nil {
		return err
	}

	if err = repo.FileWrapper.AppendLine(line); err != nil {
		return err
	}

	return nil
}

func (repo TransactionsRepo) UpdateTransaction(t TransactionDAO) error {
	line, err := repo.ToRow(t)
	if err != nil {
		return err
	}

	idxLineToUpdate := -1
	for i := 0; i < len(*repo.FileWrapper.Lines); i++ {
		if strings.Split((*repo.FileWrapper.Lines)[i], repo.FileSeparator)[0] == strconv.Itoa(t.ID) {
			idxLineToUpdate = i
		}
	}

	if err = repo.FileWrapper.ReplaceLine(idxLineToUpdate, line); err != nil {
		return err
	}

	return nil
}

func (t TransactionDAO) categoriesLabels() []string {
	var categoriesLabels []string
	for _, c := range t.Categories {
		categoriesLabels = append(categoriesLabels, c.Label)
	}
	return categoriesLabels
}
