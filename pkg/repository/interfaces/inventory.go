package interfaces

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/utils/models"
)

type InventoryRepository interface {
	AddInventory(inventory models.AddInventories) (models.InventoryResponse, error)
	ListProducts(int,int) ([]models.InventoryUserResponse, error)
	EditInventory(domain.Inventory, int) (domain.Inventory, error)
	DeleteInventory(id string) error
	CheckInventory(p_id int) (bool, error)
	UpdateInventory(p_id int, stock int) (models.InventoryResponse, error)
}
