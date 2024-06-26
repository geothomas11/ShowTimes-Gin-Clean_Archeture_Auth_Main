package usecase

import (
	"ShowTimes/pkg/domain"
	helper "ShowTimes/pkg/helper/interface"
	"fmt"

	repo "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"

	usecase "ShowTimes/pkg/usecase/interface"
)

type productUseCase struct {
	repository repo.ProductRepository
	helper     helper.Helper
}

func NewInventoryUseCase(rep repo.ProductRepository, h helper.Helper) usecase.ProductUseCase {
	return &productUseCase{
		repository: rep,
		helper:     h,
	}

}

func (i *productUseCase) AddProducts(inventory models.AddProducts) (models.ProductResponse, error) {

	InventoryResponse, err := i.repository.AddProducts(inventory)
	if err != nil {
		return models.ProductResponse{}, err
	}
	return InventoryResponse, nil

}
func (i *productUseCase) ListProducts(pageNo, pageList int) ([]models.ProductUserResponse, error) {

	offSet := (pageNo - 1) * pageList
	productList, err := i.repository.ListProducts(pageList, offSet)
	if err != nil {
		return []models.ProductUserResponse{}, nil
	}
	return productList, nil

}

func (usecase *productUseCase) EditProducts(inventory domain.Product, id int) (domain.Product, error) {
	modInventory, err := usecase.repository.EditProducts(inventory, id)
	if err != nil {
		return domain.Product{}, err
	}
	return modInventory, nil
}

func (usecase *productUseCase) DeleteProducts(inventoryID string) error {

	err := usecase.repository.DeleteProducts(inventoryID)
	if err != nil {
		fmt.Println("usecase errr", err)
		return err
	}
	return nil

}
func (i productUseCase) UpdateProducts(pid int, stock int) (models.ProductResponse, error) {

	result, err := i.repository.CheckProducts(pid)
	if err != nil {
		return models.ProductResponse{}, err
	}
	if !result {
		return models.ProductResponse{}, errors.New("there is no inventory as you mentioned")
	}
	newCat, err := i.repository.UpdateProducts(pid, stock)
	if err != nil {
		return models.ProductResponse{}, err
	}
	return newCat, err
}
