package repositories

import (
	"fmt"
models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
"github.com/h-abranches-dev/daily-expenses-be/files"
)
type Repo struct {
	FileWrapper   *files.FileWrapper
	FileSeparator string
}

var (
	transTestDebRepo      *TransactionsRepo
	transTest2DebCredRepo *TransactionsRepo
	transEntRepo          *TransactionsEntitiesRepo
	categoriesRepo        *CategoriesRepo
)

func NewRepo(dbFile string) (*Repo, error) {
	fw := files.NewFileWrapper(dbFile)
	if fw == nil {
		return nil, fmt.Errorf("database file %q not found", dbFile)
	}
	return &Repo{
		FileWrapper:   fw,
		FileSeparator: ";",
	}, nil
}

func GetAllRepos() (*[]*TransactionsRepo, error) {
	r1, err := GetTransTestDebRepo()
	if err != nil {
		return nil, err
	}
	r2, err := GetTransTest2DebCredRepo()
	if err != nil {
		return nil, err
	}
	return &[]*TransactionsRepo{r1, r2}, nil
}

func GetTransTestDebRepo() (*TransactionsRepo, error) {
	if transTestDebRepo == nil {
		var err error
		if transTestDebRepo, err = NewTransactionsRepo(string(models.DebitBankAccountKind), string(models.TestEntity)); err != nil {
			return nil, err
		}
	}
	return transTestDebRepo, nil
}

func GetTransTest2DebCredRepo() (*TransactionsRepo, error) {
	if transTest2DebCredRepo == nil {
		var err error
		if transTest2DebCredRepo, err = NewTransactionsRepo(string(models.DebitCreditBankAccountKind), string(models.Test2Entity)); err != nil {
			return nil, err
		}
	}
	return transTest2DebCredRepo, nil
}

func GetTransRepo(kind, entity string) (*TransactionsRepo, error) {
	if kind == string(models.DebitBankAccountKind) && entity == string(models.TestEntity) {
		return GetTransTestDebRepo()
	}
	if kind == string(models.DebitCreditBankAccountKind) && entity == string(models.Test2Entity) {
		return GetTransTest2DebCredRepo()
	}
	return nil, fmt.Errorf("for the kind %q and entity %q any repository was found", kind, entity)
}

func GetTransEntRepo() (*TransactionsEntitiesRepo, error) {
	if transEntRepo == nil {
		var err error
		if transEntRepo, err = NewTransactionsEntitiesRepo(); err != nil {
			return nil, err
		}
	}
	return transEntRepo, nil
}

func GetCategoriesRepo() (*CategoriesRepo, error) {
	if categoriesRepo == nil {
		var err error
		if categoriesRepo, err = NewCategoriesRepo(); err != nil {
			return nil, err
		}
	}
	return categoriesRepo, nil
}
