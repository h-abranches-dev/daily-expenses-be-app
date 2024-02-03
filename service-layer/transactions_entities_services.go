package services

import (
	"fmt"
	"strconv"
	"strings"
models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
repositories "github.com/h-abranches-dev/daily-expenses-be/persistence-layer"
)
type TransactionsEntityDTO struct {
	ID      string `json:"id"`
	Entity  string `json:"entity"`
	Kind    string `json:"type"`
	Balance string `json:"balance"`
}

type TransactionsEntitiesDTO []TransactionsEntityDTO

func newTransactionsEntityDTO(tse models.TransactionsEntity) TransactionsEntityDTO {
	return TransactionsEntityDTO{
		ID:      strconv.Itoa(tse.ID),
		Entity:  tse.Entity,
		Kind:    tse.Kind,
		Balance: fmt.Sprintf("%.2f", tse.Balance),
	}
}

func NewTransactionsEntitiesDTO(tses models.TransactionsEntities) TransactionsEntitiesDTO {
	var tsesDTO TransactionsEntitiesDTO
	for i := 0; i < len(tses); i++ {
		tsesDTO = append(tsesDTO, newTransactionsEntityDTO(tses[i]))
	}
	return tsesDTO
}

func newTransactionsEntity(tseDAO repositories.TransactionsEntityDAO) models.TransactionsEntity {
	return models.TransactionsEntity{
		ID:      tseDAO.ID,
		Entity:  tseDAO.Entity,
		Kind:    tseDAO.Kind,
		Balance: tseDAO.Balance,
	}
}

func newTransactionsEntities(tsesDAO repositories.TransactionsEntitiesDAO) models.TransactionsEntities {
	var tses models.TransactionsEntities
	for i := 0; i < len(tsesDAO); i++ {
		tses = append(tses, newTransactionsEntity(tsesDAO[i]))
	}
	return tses
}

func (tseDTO TransactionsEntityDTO) NewTransactionsEntity() (models.TransactionsEntity, error) {
	tse := models.TransactionsEntity{}
	var balance float64
	if tseDTO.Balance != "" {
		var err error
		balance, err = strconv.ParseFloat(strings.Trim(tseDTO.Balance, " "), 32)
		if err != nil {
			return tse, err
		}
	}

	tse.Entity = tseDTO.Entity
	tse.Kind = tseDTO.Kind
	tse.Balance = float32(balance)

	return tse, nil
}

func newTransactionsEntityDAO(tse models.TransactionsEntity) repositories.TransactionsEntityDAO {
	return repositories.TransactionsEntityDAO{
		ID:      tse.ID,
		Entity:  tse.Entity,
		Kind:    tse.Kind,
		Balance: tse.Balance,
	}
}

func GetAllTransactionsEntities(repo *repositories.TransactionsEntitiesRepo, tses *models.TransactionsEntities) (*models.TransactionsEntities, error) {
	if repo == nil {
		return nil, fmt.Errorf("transactions entities repo wasn't initialized")
	}
	if tses == nil {
		tsesDAO, err := repo.GetAllTransactionsEntities()
		if err != nil {
			return nil, err
		}
		tses = new(models.TransactionsEntities)
		*tses = newTransactionsEntities(tsesDAO)
	}
	return tses, nil
}

func transactionsEntitiesNextAvailableID(repo *repositories.TransactionsEntitiesRepo, tses *models.TransactionsEntities) (int, error) {
	if repo == nil {
		return -1, fmt.Errorf("transactions entities repo wasn't initialized")
	}
	if tses == nil {
		tsesDAO, err := repo.GetAllTransactionsEntities()
		if err != nil {
			return -1, err
		}
		tses = new(models.TransactionsEntities)
		*tses = newTransactionsEntities(tsesDAO)
	}
	if len(*tses) == 0 {
		return 1, nil
	}
	maxID := 0
	for i := 0; i < len(*tses); i++ {
		if (*tses)[i].ID > maxID {
			maxID = (*tses)[i].ID
		}
	}
	if maxID == 0 {
		return -1, fmt.Errorf("it wasn't found a max ID")
	}
	return maxID + 1, nil
}

func AddTransactionsEntity(repo *repositories.TransactionsEntitiesRepo, tse models.TransactionsEntity) (int, error) {
	if repo == nil {
		return -1, fmt.Errorf("transactions entities repo wasn't initialized")
	}
	tseDAO := newTransactionsEntityDAO(tse)
	var err error
	tseDAO.ID, err = transactionsEntitiesNextAvailableID(repo, models.RefToTransactionsEntities)
	if err != nil {
		return -1, err
	}
	err = repo.AddTransactionsEntity(tseDAO)
	if err != nil {
		return -1, err
	}
	return tseDAO.ID, nil
}

func updateTransactionsEntity(repo *repositories.TransactionsEntitiesRepo, tse models.TransactionsEntity) error {
	if repo == nil {
		return fmt.Errorf("transactions entities repo wasn't initialized")
	}
	tseDAO := newTransactionsEntityDAO(tse)
	err := repo.UpdateTransactionsEntity(tseDAO)
	if err != nil {
		return err
	}
	return nil
}
