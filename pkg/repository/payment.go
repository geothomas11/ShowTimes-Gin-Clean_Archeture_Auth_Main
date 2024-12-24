package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepository{
		db: db,
	}

}

func (pr *paymentRepository) PaymentExist(orderBody models.OrderIncoming) (bool, error) {
	var count int
	if err := pr.db.Raw("SELECT COUNT(*) FROM payment_methods WHERE id = ?", orderBody.PaymentID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (pr *paymentRepository) PaymentMethodID(orderID int) (int, error) {
	var a int
	err := pr.db.Raw("SELECT payment_method_id FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil

}

func (pr *paymentRepository) AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error) {
	var payment string
	if err := pr.db.Raw("INSERT INTO payment_methods (payment_name) VALUES (?) RETRUNING Payment_name ", pay.PaymentName).Scan(&payment).Error; err != nil {
		return models.PaymentDetails{}, err
	}
	var paymentResponse models.PaymentDetails
	err := pr.db.Raw("SELECT id,payment_name FROM payment_methods WHERE payment_name = ?", payment).Scan(&paymentResponse).Error
	if err != nil {
		return models.PaymentDetails{}, err
	}
	return paymentResponse, nil

}

func (pr *paymentRepository) CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
	var count int
	err := pr.db.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1", payment).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil

}
