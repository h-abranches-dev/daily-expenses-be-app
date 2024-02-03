package services

import (
	"fmt"
	models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
	repositories "github.com/h-abranches-dev/daily-expenses-be/persistence-layer"
	"strconv"
)

type CategoryDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type CategoriesDTO []CategoryDTO

func newCategoryDTO(c models.Category) CategoryDTO {
	return CategoryDTO{
		ID:    strconv.Itoa(c.ID),
		Label: c.Label,
	}
}

func NewCategoriesDTO(cs models.Categories) CategoriesDTO {
	var csDTO CategoriesDTO
	for i := 0; i < len(cs); i++ {
		csDTO = append(csDTO, newCategoryDTO(cs[i]))
	}
	return csDTO
}

func newCategory(cDAO repositories.CategoryDAO) models.Category {
	return models.Category{
		ID:    cDAO.ID,
		Label: cDAO.Label,
	}
}

func newCategories(csDAO repositories.CategoriesDAO) models.Categories {
	var cs models.Categories
	for i := 0; i < len(csDAO); i++ {
		cs = append(cs, newCategory(csDAO[i]))
	}
	return cs
}

func (cDTO CategoryDTO) NewCategory() (models.Category, error) {
	return models.Category{
		Label: cDTO.Label,
	}, nil
}

func newCategoryDAO(c models.Category) repositories.CategoryDAO {
	return repositories.CategoryDAO{
		ID:    c.ID,
		Label: c.Label,
	}
}

func GetAllCategories(repo *repositories.CategoriesRepo, cs *models.Categories) (*models.Categories, error) {
	if repo == nil {
		return nil, fmt.Errorf("categories repo wasn't initialized")
	}
	if cs == nil {
		csDAO, err := repo.GetAllCategories()
		if err != nil {
			return nil, err
		}
		cs = new(models.Categories)
		*cs = newCategories(csDAO)
	}
	return cs, nil
}

func categoriesNextAvailableID(repo *repositories.CategoriesRepo, cs *models.Categories) (int, error) {
	if repo == nil {
		return -1, fmt.Errorf("categories repo wasn't initialized")
	}
	if cs == nil {
		csDAO, err := repo.GetAllCategories()
		if err != nil {
			return -1, err
		}
		cs = new(models.Categories)
		*cs = newCategories(csDAO)
	}
	if len(*cs) == 0 {
		return 1, nil
	}
	maxID := 0
	for i := 0; i < len(*cs); i++ {
		if (*cs)[i].ID > maxID {
			maxID = (*cs)[i].ID
		}
	}
	if maxID == 0 {
		return -1, fmt.Errorf("it wasn't found a max ID")
	}
	return maxID + 1, nil
}

func AddCategory(repo *repositories.CategoriesRepo, c models.Category) (int, error) {
	if repo == nil {
		return -1, fmt.Errorf("transactions entities repo wasn't initialized")
	}
	cDAO := newCategoryDAO(c)
	var err error
	cDAO.ID, err = categoriesNextAvailableID(repo, models.RefToCategories)
	if err != nil {
		return -1, err
	}
	err = repo.AddCategory(cDAO)
	if err != nil {
		return -1, err
	}
	return cDAO.ID, nil
}

func UpdateCategory(repo *repositories.CategoriesRepo, t models.Category) error {
	if repo == nil {
		return fmt.Errorf("categories repo wasn't initialized")
	}
	cDAO := newCategoryDAO(t)
	err := repo.UpdateCategory(cDAO)
	if err != nil {
		return err
	}
	return nil
}

func filterCategories(allCategories *models.Categories, categories models.Categories) (models.Categories, error) {
	cs := models.Categories{}
	for _, cDTO := range categories {
		for _, c := range *allCategories {
			if c.Label == cDTO.Label {
				cs = append(cs, c)
				break
			}
		}
	}
	if len(categories) != 0 && len(cs) != len(categories) {
		return models.Categories{}, fmt.Errorf("the categories %v don't match", categories)
	}
	return cs, nil
}
