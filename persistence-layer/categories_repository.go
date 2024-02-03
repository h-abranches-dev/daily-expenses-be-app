package repositories

import (
	"fmt"
	"strconv"
	"strings"
)

type CategoryDAO struct {
	ID    int
	Label string
}

type CategoriesDAO []CategoryDAO

type CategoriesRepo struct {
	*Repo
}

const (
	categoriesDBFile string = "db/categories.csv"
)

func NewCategoriesRepo() (*CategoriesRepo, error) {
	r, err := NewRepo(categoriesDBFile)
	if err != nil {
		return nil, err
	}
	return &CategoriesRepo{
		Repo: r,
	}, nil
}

func (repo CategoriesRepo) ToRow(c CategoryDAO) (string, error) {
	if strings.Index(c.Label, repo.FileSeparator) != -1 {
		return "", fmt.Errorf("invalid 'Label' because includes the char %q => %q", repo.FileSeparator,
			c.Label)
	}

	return fmt.Sprintf("%d;%s", c.ID, c.Label), nil
}

func (repo CategoriesRepo) rowToCategory(row string) (CategoryDAO, error) {
	emptyCategory := CategoryDAO{}
	columns := strings.Split(row, ";")
	cID, err := strconv.Atoi(columns[0])
	if err != nil {
		return emptyCategory, err
	}

	return CategoryDAO{
		ID:    cID,
		Label: columns[1],
	}, nil
}

func (repo CategoriesRepo) GetAllCategories() (CategoriesDAO, error) {
	categories := CategoriesDAO{}
	rows := repo.FileWrapper.Lines

	for i := 1; i < len(*rows)-1; i++ {
		category, err := repo.rowToCategory((*rows)[i])
		if err != nil {
			return CategoriesDAO{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (repo CategoriesRepo) AddCategory(c CategoryDAO) error {
	line, err := repo.ToRow(c)
	if err != nil {
		return err
	}

	if err = repo.FileWrapper.AppendLine(line); err != nil {
		return err
	}

	return nil
}

func (repo CategoriesRepo) UpdateCategory(c CategoryDAO) error {
	line, err := repo.ToRow(c)
	if err != nil {
		return err
	}

	idxLineToUpdate := -1
	for i := 0; i < len(*repo.FileWrapper.Lines); i++ {
		if strings.Split((*repo.FileWrapper.Lines)[i], repo.FileSeparator)[0] == strconv.Itoa(c.ID) {
			idxLineToUpdate = i
		}
	}

	if err = repo.FileWrapper.ReplaceLine(idxLineToUpdate, line); err != nil {
		return err
	}

	return nil
}

func (repo CategoriesRepo) CategoryDAO(categoryLabel string) (CategoryDAO, error) {
	emptyCategory := CategoryDAO{}
	cs, err := repo.GetAllCategories()
	if err != nil {
		return emptyCategory, err
	}
	for _, c := range cs {
		if c.Label == categoryLabel {
			return c, nil
		}
	}
	return emptyCategory, fmt.Errorf("category %q not found", categoryLabel)
}

func (repo CategoriesRepo) CategoriesDAO(categoriesLabels []string) (CategoriesDAO, error) {
	emptyCategories := CategoriesDAO{}
	csDAO := make(CategoriesDAO, 0)
	for _, l := range categoriesLabels {
		c, err := repo.CategoryDAO(l)
		if err != nil {
			return emptyCategories, err
		}
		csDAO = append(csDAO, c)
	}
	if len(csDAO) == 0 {
		return CategoriesDAO{}, fmt.Errorf("it wasn´t found any supported category")
	}
	return csDAO, nil
}
