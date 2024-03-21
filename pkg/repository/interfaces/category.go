package interfaces

import "ShowTimes/pkg/domain"

type CategoryRepository interface {
	AddCatogery(category domain.Category) (domain.Category, error)
	GetCategories() ([]domain.Category, error)
	UpdateCategory(current, new string) (domain.Category, error)
	CheckCategory(current string) (bool, error)
	DeleteCategory(categoryID int) error
	IsCategoryExist(category string) (bool, error)
}
