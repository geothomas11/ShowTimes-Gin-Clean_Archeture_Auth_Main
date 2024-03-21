package usecase

import (
	"ShowTimes/pkg/domain"
	repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"strconv"

	"errors"
)

type categoryUseCase struct {
	repository repo.CategoryRepository
}

func NewCategoryUseCase(repo repo.CategoryRepository) interfaces.CategoryUseCase {
	return &categoryUseCase{
		repository: repo,
	}

}

func (cat *categoryUseCase) AddCategory(category domain.Category) (domain.Category, error) {
	if category.Category == "" {
		return domain.Category{}, errors.New("category name not valid")
	}
	exist, err := cat.repository.IsCategoryExist(category.Category)
	if err != nil {
		return domain.Category{}, err
	}
	if exist {
		return domain.Category{}, errors.New("category already exist")
	}
	productResponse, err := cat.repository.AddCatogery(category)
	if err != nil {
		return domain.Category{}, err
	}
	return productResponse, nil
}

func (cat *categoryUseCase) GetCategories() ([]domain.Category, error) {
	categories, err := cat.repository.GetCategories()
	if err != nil {
		return []domain.Category{}, err
	}
	return categories, nil

}
func (cat *categoryUseCase) UpdateCategory(current string, new string) (domain.Category, error) {
	result, err := cat.repository.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}
	if !result {
		return domain.Category{}, errors.New("there is no category as you mentioned")
	}
	newcat, err := cat.repository.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}
	return newcat, err
}

func (cat *categoryUseCase) DeleteCategory(categoryID string) error {
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return errors.New("string conversion invalid")
	}
	err = cat.repository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil

}
