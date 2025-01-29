package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"database/sql"
	"errors"
	"fmt"

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
	if err := pr.db.Exec("INSERT INTO payment_methods (payment_name) VALUES (?)  ", pay.PaymentName).Error; err != nil {
		return models.PaymentDetails{}, err
	}
	var paymentResponse models.PaymentDetails
	err := pr.db.Raw("SELECT id,payment_name FROM payment_methods WHERE payment_name = ?", pay.PaymentName).Scan(&paymentResponse).Error
	fmt.Println("error in repo", err)
	if err != nil {
		if err == sql.ErrNoRows {

			return models.PaymentDetails{}, errors.New("no data found")

		}
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

//RAZOR PAY INTEGRATION

func (repo *ProductRepository) AddRazorPayDetails(orderId int, razorPayId string) error {
	query := `INSERT INTO payments (order_id,razorpayid) values($1,$2) `
	if err := repo.DB.Exec(query, orderId, razorPayId).Error; err != nil {
		err = errors.New("error in inserting values to razor pay data table" + err.Error())
		return err
	}
	return nil

}
func (pr *paymentRepository) UpdatePaymentDetails(orderId string, paymentId string) error {

	if err := pr.db.Exec("update payments set payment = $1 where razer_id = $2", paymentId, orderId).Error; err != nil {
		err = errors.New("error in updating the razer pay table " + err.Error())
		return err
	}
	return nil
}

func (pr *paymentRepository) GetPaymentStatus(orderId string) (bool, error) {
	var paymentStatus string
	err := pr.db.Raw("SELECT  payment_status from orders WHERE id = $1", orderId).Scan(&paymentStatus).Error
	if err != nil {
		return false, err
	}

	// Check if payment status is "PAID"
	isPaid := paymentStatus == "PAID"

	return isPaid, nil
}

func (pr *paymentRepository) UpdatePaymentStatus(status bool, orderId string) error {
	var paymentStatus string
	if status {
		paymentStatus = "PAID"
	} else {
		paymentStatus = "NOT PAID"
	}
	query := `UPDATE orders SET payment_status = 'SHIPPED' WHERE id = $2`
	if err := pr.db.Exec(query, paymentStatus, orderId).Error; err != nil {
		err = errors.New("error in updating orders payment status" + err.Error())
		return err
	}
	return nil

}
