package interfaces

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/utils/models"
)

type InventoryUseCase interface {
	AddInventory(inventory models.AddInventories) (models.InventoryResponse, error)
	ListProducts(int, int) ([]models.InventoryUserResponse, error)
	EditInventory(domain.Inventory, int) (domain.Inventory, error)
	DeleteInventory(id string) error
	UpdateInventory(poductID int, stock int) (models.InventoryResponse, error)
}
