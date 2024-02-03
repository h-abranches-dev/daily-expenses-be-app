package repositories

import (
	"fmt"
	"strconv"
	"strings"
)

type TransactionsEntityDAO struct {
	ID      int
	Entity  string
	Kind    string
	Balance float32
}

type TransactionsEntitiesDAO []TransactionsEntityDAO

type TransactionsEntitiesRepo struct {
	*Repo
}

const (
	transactionsEntitiesDBFile string = "db/transactions_entities.csv"
)

func NewTransactionsEntitiesRepo() (*TransactionsEntitiesRepo, error) {
	r, err := NewRepo(transactionsEntitiesDBFile)
	if err != nil {
		return nil, err
	}
	return &TransactionsEntitiesRepo{
		Repo: r,
	}, nil
}

func rowToTransactionsEntity(row string) (TransactionsEntityDAO, error) {
	emptyTransactionsEntity := TransactionsEntityDAO{}
	columns := strings.Split(row, ";")
	id, err := strconv.Atoi(columns[0])
	if err != nil {
		return emptyTransactionsEntity, err
	}
	balance, err := strconv.ParseFloat(strings.Trim(columns[3], " "), 32)
	if err != nil {
		return emptyTransactionsEntity, err
	}
	return TransactionsEntityDAO{
		ID:      id,
		Entity:  columns[1],
		Kind:    columns[2],
		Balance: float32(balance),
	}, nil
}

func (repo TransactionsEntitiesRepo) GetTransactionsEntity(entity, kind string) (TransactionsEntityDAO, error) {
	tseDAO := TransactionsEntityDAO{}
	tsesDAO, err := repo.GetAllTransactionsEntities()
	if err != nil {
		return tseDAO, err
	}
	for _, v := range tsesDAO {
		if v.Entity == entity && v.Kind == kind {
			tseDAO = v
			break
		}
	}
	return tseDAO, nil
}

func (repo TransactionsEntitiesRepo) GetAllTransactionsEntities() (TransactionsEntitiesDAO, error) {
	transactionsEntities := TransactionsEntitiesDAO{}
	rows := repo.FileWrapper.Lines

	for i := 1; i < len(*rows)-1; i++ {
		transactionsEntity, err := rowToTransactionsEntity((*rows)[i])
		if err != nil {
			return TransactionsEntitiesDAO{}, err
		}
		transactionsEntities = append(transactionsEntities, transactionsEntity)
	}
	return transactionsEntities, nil
}

func (repo TransactionsEntitiesRepo) ToRow(tse TransactionsEntityDAO) (string, error) {
	return fmt.Sprintf("%d;%s;%s;%s", tse.ID, tse.Entity,
		tse.Kind, fmt.Sprintf("%.2f", tse.Balance)), nil
}

func (repo TransactionsEntitiesRepo) AddTransactionsEntity(tse TransactionsEntityDAO) error {
	line, err := repo.ToRow(tse)
	if err != nil {
		return err
	}

	if err = repo.FileWrapper.AppendLine(line); err != nil {
		return err
	}

	return nil
}

func (repo TransactionsEntitiesRepo) UpdateTransactionsEntity(tse TransactionsEntityDAO) error {
	line, err := repo.ToRow(tse)
	if err != nil {
		return err
	}

	idxLineToUpdate := -1
	for i := 0; i < len(*repo.FileWrapper.Lines); i++ {
		if strings.Split((*repo.FileWrapper.Lines)[i], repo.FileSeparator)[0] == strconv.Itoa(tse.ID) {
			idxLineToUpdate = i
		}
	}

	if idxLineToUpdate == -1 {
		return fmt.Errorf("line to update wasn't found")
	}

	if err = repo.FileWrapper.ReplaceLine(idxLineToUpdate, line); err != nil {
		return err
	}

	return nil
}
