package repository

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type OfferRepository struct {
	DB *gorm.DB
}

func NewOfferRepository(DB *gorm.DB) interfaces.OfferRepository {
	return &OfferRepository{DB: DB}

}
func (or *OfferRepository) AddProductOffer(ProductOffer models.ProductOfferResp) error {
	var count int
	err := or.DB.Raw("SELECT count(*) from product_offers where offer_name = ? and product_id = ?", ProductOffer.OfferName, ProductOffer.ProductID).Scan(&count).Error
	if err != nil {
		return errors.New("error in getting offer details")
	}
	if count > 0 {
		return errors.New("offer already exist")
	}
	err = or.DB.Raw("select count(*) from product_offers where product_id = ?", ProductOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = or.DB.Exec("DELETE FROM product_offers WHERE product_id = ?", ProductOffer.ProductID).Scan(&count).Error
		if err != nil {
			return err
		}
	}
	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = or.DB.Exec("INSERT INTO product_offers (product_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)", ProductOffer.ProductID, ProductOffer.OfferName, ProductOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}
	return nil

}

func (or *OfferRepository) AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error {
	var count int

	// ✅ Check if the category exists before adding an offer
	var categoryCount int
	err := or.DB.Raw("SELECT COUNT(*) FROM categories WHERE id = ?", CategoryOffer.CategoryID).Scan(&categoryCount).Error
	if err != nil {
		return errors.New("error checking category existence")
	}
	if categoryCount == 0 {
		return errors.New("invalid category_id: the referenced category does not exist")
	}

	// ✅ Check if the same offer already exists for the category
	err = or.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE offer_name = ? AND category_id = ?", CategoryOffer.OfferName, CategoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return errors.New("error in getting offer details")
	}
	if count > 0 {
		return errors.New("offer already exists")
	}

	// ✅ Check if another offer exists for the same category and delete it
	err = or.DB.Raw("SELECT COUNT(*) FROM category_offers WHERE category_id = ?", CategoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err = or.DB.Exec("DELETE FROM category_offers WHERE category_id = ?", CategoryOffer.CategoryID).Error
		if err != nil {
			return err
		}
	}
	startDate := time.Now()
	endDate := startDate.Add(time.Hour * 24 * 5)

	err = or.DB.Exec("INSERT INTO category_offers (category_id, offer_name, discount_percentage, start_date, end_date) VALUES (?, ?, ?, ?, ?)",
		CategoryOffer.CategoryID, CategoryOffer.OfferName, CategoryOffer.DiscountPercentage, startDate, endDate).Error
	if err != nil {
		return err
	}

	return nil
}

func (or *OfferRepository) GetProductOffer() ([]domain.ProductOffer, error) {
	var producttOfferDetails []domain.ProductOffer
	err := or.DB.Raw("SELECT * from product_offers").Scan(&producttOfferDetails).Error
	if err != nil {
		return []domain.ProductOffer{}, errors.New("error in getting category offers")
	}
	return producttOfferDetails, nil

}
func (or *OfferRepository) GetCategoryOffer() ([]domain.CategoryOffer, error) {
	var CategorytOfferDetails []domain.CategoryOffer
	err := or.DB.Raw("SELECT * from category_offers").Scan(&CategorytOfferDetails).Error
	if err != nil {
		return []domain.CategoryOffer{}, errors.New("error in getting category offers")
	}
	return CategorytOfferDetails, nil

}
func (or *OfferRepository) ExpireProductOffer(id int) error {
	if err := or.DB.Exec("DELETE FROM product_offers WHERE id = $1", id).Error; err != nil {
		return err
	}
	return nil

}
func (or *OfferRepository) ExpireCategoryOffer(id int) error {
	if err := or.DB.Exec("DELETE FROM category_offers WHERE id = $1", id).Error; err != nil {
		return err
	}
	return nil

}

func (or *OfferRepository) GetCatOfferPercent(categoryId int) (int, error) {
	var percent int
	err := or.DB.Raw("select discount_percentage from category_offers where category_id = ?", categoryId).Scan(&percent).Error
	if err != nil {
		return 0, errors.New(errmsg.ErrGetDB)
	}
	return percent, nil

}

func (or *OfferRepository) GetProOfferPercent(productId int) (int, error) {
	var percent int
	err := or.DB.Raw("select discount_percentage from category_offers where product_id = ?", productId).Scan(&percent).Error
	if err != nil {
		return 0, errors.New(errmsg.ErrGetDB)
	}
	return percent, nil

}
