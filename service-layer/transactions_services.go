package services

import (
	"fmt"
	"strconv"
	"strings"
models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
repositories "github.com/h-abranches-dev/daily-expenses-be/persistence-layer"
"github.com/h-abranches-dev/daily-expenses-be/utils"
"time"
)
type TransactionDTO struct {
	ID              string   `json:"id"`
	TransactionDate string   `json:"transaction_date"`
	Transaction     string   `json:"transaction"`
	Categories      []string `json:"categories"`
	Kind            string   `json:"type,omitempty"`
	Amount          float32  `json:"amount"`
}

type TransactionsDTO struct {
	TransactionsDTO []TransactionDTO `json:"transactions"`
	Total           *float32         `json:"total,omitempty"`
}

func newTransactionDTO(t models.Transaction) (TransactionDTO, error) {
	categoriesRepo, err := repositories.GetCategoriesRepo()
	if err != nil {
		return TransactionDTO{}, err
	}

	categoriesDTO, err := GetAllCategories(categoriesRepo, models.RefToCategories)
	if err != nil {
		return TransactionDTO{}, err
	}

	*categoriesDTO, err = filterCategories(categoriesDTO, t.Categories)
	if err != nil {
		return TransactionDTO{}, err
	}

	var categoriesLabels []string
	for _, c := range *categoriesDTO {
		categoriesLabels = append(categoriesLabels, c.Label)
	}

	return TransactionDTO{
		ID:              strconv.Itoa(t.ID),
		TransactionDate: t.TransactionDate.Format(utils.DateFormat),
		Transaction:     t.Transaction,
		Categories:      categoriesLabels,
		Kind:            t.Kind,
		Amount:          t.Amount,
	}, nil
}

func NewTransactionsDTO(ts models.Transactions) (TransactionsDTO, error) {
	var tsDTO TransactionsDTO
	for i := 0; i < len(ts); i++ {
		ntDTO, err := newTransactionDTO(ts[i])
		if err != nil {
			return TransactionsDTO{}, err
		}
		tsDTO.TransactionsDTO = append(tsDTO.TransactionsDTO, ntDTO)
	}
	return tsDTO, nil
}

func newTransaction(tDAO repositories.TransactionDAO) models.Transaction {
	categories := make([]models.Category, 0)
	for _, c := range tDAO.Categories {
		categories = append(categories, models.Category(c))
	}

	return models.Transaction{
		ID:              tDAO.ID,
		TransactionDate: tDAO.TransactionDate,
		Transaction:     tDAO.Transaction,
		Categories:      categories,
		Kind:            tDAO.Kind,
		Amount:          tDAO.Amount,
	}
}

func newTransactions(tsDAO repositories.TransactionsDAO) models.Transactions {
	var ts models.Transactions
	for i := 0; i < len(tsDAO); i++ {
		ts = append(ts, newTransaction(tsDAO[i]))
	}
	return ts
}

func NewCategory(categoryLabel string) (models.Category, error) {
	return models.Category{
		Label: categoryLabel,
	}, nil
}

func (tDTO TransactionDTO) NewTransaction() (models.Transaction, error) {
	t := models.Transaction{}
	if (tDTO.Kind != models.DebitKindTransaction) && (tDTO.Kind != models.CreditKindTransaction) && (tDTO.Kind != "") {
		return t, fmt.Errorf("the value %q for type field is not valid", tDTO.Kind)
	}
	tDate, err := time.Parse(utils.DateFormat, strings.Trim(tDTO.TransactionDate, " "))
	if err != nil {
		return t, err
	}

	categories := make([]models.Category, 0)
	for _, cl := range tDTO.Categories {
		var nc models.Category
		nc, err = NewCategory(cl)
		if err != nil {
			return t, err
		}
		categories = append(categories, nc)
	}

	t.TransactionDate = tDate
	t.Transaction = strings.Trim(tDTO.Transaction, " ")
	t.Categories = categories
	t.Kind = tDTO.Kind
	t.Amount = tDTO.Amount

	return t, nil
}

func newTransactionDAO(t models.Transaction) repositories.TransactionDAO {
	categories := make(repositories.CategoriesDAO, 0)
	for _, c := range t.Categories {
		categories = append(categories, newCategoryDAO(c))
	}

	return repositories.TransactionDAO{
		ID:              t.ID,
		TransactionDate: t.TransactionDate,
		Transaction:     t.Transaction,
		Categories:      categories,
		Kind:            t.Kind,
		Amount:          t.Amount,
	}
}

func GetAllTransactionsByRepo(repo *repositories.TransactionsRepo, mandatoryUseOfDB bool) (*models.Transactions, error) {
	if repo == nil {
		return nil, fmt.Errorf("transactions repo wasn't initialized")
	}
	key := models.EntityKindKey(models.TransactionEntity(repo.Entity), models.TransactionKind(repo.Kind))
	refToTransactions := (*models.RefsToTransactions)[key]
	if refToTransactions == nil || mandatoryUseOfDB {
		tsDAO, err := repo.GetAllTransactions()
		if err != nil {
			return nil, err
		}
		refToTransactions = new(models.Transactions)
		*refToTransactions = newTransactions(tsDAO)
	}
	(*models.RefsToTransactions)[key] = refToTransactions
	return refToTransactions, nil
}

func GetAllTransactions(repos []*repositories.TransactionsRepo, mandatoryUseOfDB bool) (*models.Transactions, error) {
	var allTransactions = new(models.Transactions)
	for _, r := range repos {
		if r == nil {
			return nil, fmt.Errorf("transactions repo wasn't initialized")
		}
		key := models.EntityKindKey(models.TransactionEntity(r.Entity), models.TransactionKind(r.Kind))
		refToTransactions := (*models.RefsToTransactions)[key]
		if refToTransactions == nil || mandatoryUseOfDB {
			tsDAO, err := r.GetAllTransactions()
			if err != nil {
				return nil, err
			}
			refToTransactions = new(models.Transactions)
			*refToTransactions = newTransactions(tsDAO)
		}
		(*models.RefsToTransactions)[key] = refToTransactions

		*allTransactions = append(*allTransactions, *refToTransactions...)
	}

	return allTransactions, nil
}

func transactionsNextAvailableID(repo *repositories.TransactionsRepo) (int, error) {
	if repo == nil {
		return -1, fmt.Errorf("transactions repo wasn't initialized")
	}
	refToTransactions := (*models.RefsToTransactions)[fmt.Sprintf("%s_#_%s", repo.Entity, repo.Kind)]
	if refToTransactions == nil {
		tsDAO, err := repo.GetAllTransactions()
		if err != nil {
			return -1, err
		}
		refToTransactions = new(models.Transactions)
		*refToTransactions = newTransactions(tsDAO)
	}
	if len(*refToTransactions) == 0 {
		return 1, nil
	}
	maxID := 0
	for i := 0; i < len(*refToTransactions); i++ {
		if (*refToTransactions)[i].ID > maxID {
			maxID = (*refToTransactions)[i].ID
		}
	}
	if maxID == 0 {
		return -1, fmt.Errorf("it wasn't found a max ID")
	}
	return maxID + 1, nil
}

func getCurrentBalance(repo *repositories.TransactionsRepo) (float32, error) {
	sum := float32(0.0)
	ts, err := GetAllTransactionsByRepo(repo, false)
	if err != nil {
		return 0.0, err
	}
	for _, t := range *ts {
		sum += t.Amount
	}
	return sum, nil
}

func updateRefToTransactions(repo *repositories.TransactionsRepo) error {
	if repo == nil {
		return fmt.Errorf("transactions repo wasn't initialized")
	}
	if _, err := GetAllTransactionsByRepo(repo, true); err != nil {
		return err
	}
	return nil
}

func AddTransaction(repo *repositories.TransactionsRepo, t models.Transaction) (int, error) {

	if repo == nil {
		return -1, fmt.Errorf("transactions repo wasn't initialized")
	}

	var err error
	tDAO := newTransactionDAO(t)
	if tDAO.ID, err = transactionsNextAvailableID(repo); err != nil {
		return -1, err
	}

	if err = repo.AddTransaction(tDAO); err != nil {
		return -1, err
	}

	if err = updateRefToTransactions(repo); err != nil {
		return -1, err
	}

	tsesRepo, err := repositories.NewTransactionsEntitiesRepo()
	if err != nil {
		return -1, err
	}

	err = updateBalance(tsesRepo, repo)
	if err != nil {
		return -1, err
	}

	return tDAO.ID, nil
}

func UpdateTransaction(repo *repositories.TransactionsRepo, t models.Transaction) error {
	if repo == nil {
		return fmt.Errorf("transactions repo wasn't initialized")
	}

	var err error
	tDAO := newTransactionDAO(t)
	if err = repo.UpdateTransaction(tDAO); err != nil {
		return err
	}

	if err = updateRefToTransactions(repo); err != nil {
		return err
	}

	tsesRepo, err := repositories.NewTransactionsEntitiesRepo()
	if err != nil {
		return err
	}
	err = updateBalance(tsesRepo, repo)
	if err != nil {
		return err
	}

	return nil
}

func updateBalance(tsesRepo *repositories.TransactionsEntitiesRepo, tsRepo *repositories.TransactionsRepo) error {

	if tsesRepo == nil {
		return fmt.Errorf("transactions entities repo wasn't initialized")
	}
	if tsRepo == nil {
		return fmt.Errorf("transactions repo wasn't initialized")
	}

	tseDAO, err := tsesRepo.GetTransactionsEntity(tsRepo.Entity, tsRepo.Kind)
	if err != nil {
		return err
	}

	balance, err := getCurrentBalance(tsRepo)
	if err != nil {
		return err
	}

	tseDAO.Balance = balance

	if err = tsesRepo.UpdateTransactionsEntity(tseDAO); err != nil {
		return err
	}

	return nil
}

func GetCategoriesLabelsInTransactionsMap(transactions *models.Transactions) map[string]int {
	m := make(map[string]int)
	for _, t := range *transactions {
		for _, cl := range t.Categories {
			m[cl.Label]++
		}
	}

	return m
}

func GetTransactionsFilteredByCategories(transactions *models.Transactions, categoriesLabelsProvided []string) (*TransactionsDTO, error) {
	tsDTO := &TransactionsDTO{}
	total := new(float32)
	*total = float32(0.0)
	for _, t := range *transactions {
		match := true
		counter := 0
		match = false
		for _, clp := range categoriesLabelsProvided {
			for _, c := range t.Categories {
				if clp == c.Label {
					counter++
					if counter == len(categoriesLabelsProvided) {
						match = true
						break
					}
				}
			}
			if match {
				break
			}
		}
		if match {
			ntDTO, err := newTransactionDTO(t)
			if err != nil {
				return nil, err
			}
			tsDTO.TransactionsDTO = append(tsDTO.TransactionsDTO, ntDTO)
			*total += t.Amount
		}
		tsDTO.Total = total
	}

	return tsDTO, nil
}
