package usecase

import (
	"ShowTimes/pkg/domain"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"

	"github.com/pkg/errors"
)

type offerUseCase struct {
	repo interfaces_repo.OfferRepository
}

func NewOfferUsecase(repo interfaces_repo.OfferRepository) interfaces.OfferUsecase {
	return &offerUseCase{
		repo: repo,
	}

}
func (ou *offerUseCase) AddProductOffer(ProductOffer models.ProductOfferResp) error {
	if err := ou.repo.AddProductOffer(ProductOffer); err != nil {
		return errors.New("error in adding product offer")
	}
	return nil

}
func (ou *offerUseCase) AddCategoryOffer(CategoryOffer models.CategorytOfferResp) error {
	if err := ou.repo.AddCategoryOffer(CategoryOffer); err != nil {
		return errors.New("error in adding category offer")
	}
	return nil

}
func (ou *offerUseCase) GetProductOffer() ([]domain.ProductOffer, error) {
	offers, err := ou.repo.GetProductOffer()
	if err != nil {
		return []domain.ProductOffer{}, err
	}
	return offers, nil

}
func (ou *offerUseCase) GetCategoryOffer() ([]domain.CategoryOffer, error) {
	offer, err := ou.repo.GetCategoryOffer()
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return offer, nil
}
func (ou *offerUseCase) ExpireProductOffer(id int) error {
	if err := ou.repo.ExpireProductOffer(id); err != nil {
		return err
	}
	return nil

}
func (ou *offerUseCase) ExpireCategoryOffer(id int) error {
	if err := ou.repo.ExpireCategoryOffer(id); err != nil {
		return err
	}
	return nil
}
