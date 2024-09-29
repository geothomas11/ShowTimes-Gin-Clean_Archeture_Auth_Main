package repository

import (
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(Db *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		db: Db,
	}

}

func (or *orderRepository) GetAllPaymentOption() ([]models.PaymentDetails, error) {
	var paymentMethods []models.PaymentDetails
	err := or.db.Raw("SELECT * FROM payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}
	return paymentMethods, nil

}

func (or *orderRepository) GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	var addressId int
	if err := or.db.Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := or.db.Raw("SELECT * FROM address WHERE id=?", addressId).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second in address")
	}
	return addressInfoResponse, nil

}

func (or *orderRepository) GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := or.db.Raw("DELECT id,final_price,shipment_status,payment_status FROM orders WHERE id = ?", orderID).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil

}

func (or *orderRepository) GetProductsInCart(cart_id int) ([]int, error) {
	var cart_products []int
	if err := or.db.Raw("SELECT product_id FROM cart_items WHERE cart_id= ?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}
	return cart_products, nil

}
func (or *orderRepository) FindProductNames(product_id int) (string, error) {
	var product_name string

	if err := or.db.Raw("SELECT name FROM products WHERE id = ?", product_id).Scan(&product_name).Error; err != nil {
		return "", err
	}
	return product_name, nil

}

func (or *orderRepository) FindCartQuantity(cart_id, product_id int) (int, error) {
	var quantity int
	if err := or.db.Raw("SELECT quantity FROM cart_items WHERE cart_id = $1 and product_id = $2", cart_id, product_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}
	return quantity, nil

}
func (or *orderRepository) FindPrice(product_id int) (float64, error) {

	var price float64
	if err := or.db.Raw("SELECT price FROM products WHERE id =? ", product_id).Scan(&price).Error; err != nil {
		return 0, nil
	}
	return price, nil
}

func (or *orderRepository) FindStock(id int) (int, error) {
	var stock int
	err := or.db.Raw("SELECT stock FROM products WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}
	return id, nil

}
