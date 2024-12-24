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

func (or *orderRepository) CheckOrderID(orderId int) (bool, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (or *orderRepository) OrderExist(orderID int) error {
	err := or.db.Raw("SELECT id FROM orders WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return err

}
func (or *orderRepository) GetShipmentStatus(orderID int) (string, error) {
	var status string
	err := or.db.Raw("SELECT shipment_status FROM orders WHERE id = ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func (or *orderRepository) UpdateOrder(orderID int) error {
	err := or.db.Exec("UPDATE orders SET S shipment_status ='processing' WHERE id =?", orderID).Error
	if err != nil {
		return err
	}
	return nil

}

func (or *orderRepository) AddOrderProducts(order_id int, cart []models.Cart) error {
	query := `
	INSERT INTO order_items (order_id,product_id,quantity,total_price)
	VALUES (?,?,?,?)`
	for _, v := range cart {
		var productID int
		if err := or.db.Raw("SELECT id from products WHERE product_name = $1", v.ProductName).Scan(&productID).Error; err != nil {
			return err
		}
		if err := or.db.Exec(query, order_id, productID, v.Quantity, v.TotalPrice).Error; err != nil {
			return err
		}
	}
	return nil

}

func (or *orderRepository) GetBriefOrderDetails(orderID int) (models.OrderSuccessResponse, error) {
	var OrderSuccessResponse models.OrderSuccessResponse
	err := or.db.Raw(`SELECT id as order_id,shipment_status FROM orders WHERE id = ?`, orderID).Scan(&OrderSuccessResponse).Error
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	return OrderSuccessResponse, nil

}
