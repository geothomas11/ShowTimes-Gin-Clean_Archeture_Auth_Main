package repository

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/repository/interfaces"
	"errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &categoryRepository{db}
}

func (cr *categoryRepository) AddCatogery(c domain.Category) (domain.Category, error) {
	var b domain.Category
	err := cr.db.Raw("INSERT INTO categories (category) VALUES (?) RETURNING*", c.Category).Scan(&b).Error
	if err != nil {
		return domain.Category{}, err
	}
	return b, nil
}
func (cr *categoryRepository) GetCategories() ([]domain.Category, error) {
	var Model []domain.Category
	err := cr.db.Raw("SELECT * FROM categories").Scan(&Model).Error
	if err != nil {
		return []domain.Category{}, err
	}
	return Model, nil
}

func (cr *categoryRepository) CheckCategory(current string) (bool, error) {
	var i int
	err := cr.db.Raw("SELECT COUNT(*) FROM categories WHERE category=?", current).Scan(&i).Error
	if err != nil {
		return false, err
	}
	if i == 0 {
		return false, err
	}
	return true, err
}

func (cr *categoryRepository) UpdateCategory(current, new string) (domain.Category, error) {

	//CHECK DATABASE CONNECTION
	if cr.db == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}
	//Update the category
	if err := cr.db.Exec("UPDATE categories SET category = $1 WHERE category = $2", new, current).Error; err != nil {
		return domain.Category{}, err
	}
	//Retrive the updated category
	var newCat domain.Category
	if err := cr.db.First(&newCat, "category=?", new).Error; err != nil {
		return domain.Category{}, err
	}
	return newCat, nil
}

func (cr *categoryRepository) DeleteCategory(categoryID int) error {

	result := cr.db.Exec("DELETE FROM categories where id=?", categoryID)

	if result.RowsAffected < 1 {
		return errors.New("no rows with that id exists")
	}
	return nil
}

// CHECKING AND ERROR HANDLING
func (cr *categoryRepository) IsCategoryExist(category string) (bool, error) {
	var count int
	if err := cr.db.Raw("SELECT COUNT(*) FROM categories where category=?", category).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
