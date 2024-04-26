package repository

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	DB *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) interfaces.InventoryRepository {
	return &inventoryRepository{
		DB: db,
	}

}
func (i *inventoryRepository) AddInventory(inventory models.AddInventories) (models.InventoryResponse, error) {
	var count int64
	i.DB.Model(&models.Inventory{}).Where("product_name=? AND category_id =?", inventory.ProductName, inventory.CategoryID).Count(&count)
	if count > 0 {
		return models.InventoryResponse{}, errors.New("product already exist in the database")
	}
	if inventory.Stock < 0 || inventory.Price < 0 {
		return models.InventoryResponse{}, errors.New("stock and price cannot be negetive")
	}
	query := `
	INSERT INTO inventories (category_id, product_name,color,stock,price)
	VALUES(?,?,?,?)
	`
	err := i.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price).Error
	if err != nil {
		return models.InventoryResponse{}, err
	}
	var inventoryRepository models.InventoryResponse
	return inventoryRepository, nil
}

func (prod *inventoryRepository) ListProducts(pageList, offset int) ([]models.InventoryUserResponse, error) {
	var product_list []models.InventoryUserResponse

	query := "SELECT i.id,i.category_id,c.category,i.product_name,i.color,i.price FROM i INNER JOIN categories c ON i.category_id = c.id LIMIT $1 OFFSET $2"
	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.InventoryUserResponse{}, errors.New("error checking the product details")
	}
	return product_list, nil

}

func (db *inventoryRepository) EditInventory(inventory domain.Inventory, id int) (domain.Inventory, error) {
	var modInventory domain.Inventory

	query := "UPDATE inventories SET categories_id =?,product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price, id).Error; err != nil {
		return domain.Inventory{}, err
	}
	if err := db.DB.First(&modInventory, id).Error; err != nil {
		return domain.Inventory{}, err
	}
	return modInventory, nil
}

func (i *inventoryRepository) DeleteInventory(inventoryID string) error {
	id, err := strconv.Atoi(inventoryID)
	if err != nil {
		return errors.New("converting into interger is not happened")
	}
	result := i.DB.Raw("DELETE FROM inventories WHERE id=?", id)
	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}
	return nil
}
func (i *inventoryRepository) CheckInventory(pid int) (bool, error) {
	var k int
	err := i.DB.Raw("SELECTV COUNT(*) FROM inventories WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}
	if k == 0 {
		return false, err
	}
	return true, err
}

func (i *inventoryRepository) UpdateInventory(pid int, stock int) (models.InventoryResponse, error) {
	//Check the Database Connection
	if i.DB == nil {
		return models.InventoryResponse{}, errors.New("Database connection is nil")
	}

	//Update the stock
	if err := i.DB.Exec("UPDATE inventories SET stock +$1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.InventoryResponse{}, err

	}

	//Retrive the update
	var newDetails models.InventoryResponse
	var newStock int
	if err := i.DB.Raw("SELECT stock FROM inventories WHERE id =? ", pid).Scan(&newStock).Error; err != nil {
		return models.InventoryResponse{}, err
	}
	newDetails.ProductID = pid
	newDetails.Stock = newStock
	return newDetails, nil
}
