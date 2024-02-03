package models

type Category struct {
	ID    int
	Label string
}

type Categories []Category

var (
	RefToCategories *Categories
)
