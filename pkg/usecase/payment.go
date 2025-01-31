package usecase

import (
	"ShowTimes/pkg/config"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"errors"

	"github.com/razorpay/razorpay-go"
)

type paymentUseCase struct {
	paymentRepository interfaces_repo.PaymentRepository
	order_Repo        interfaces_repo.OrderRepository
	cfg               config.Config
}

func NewPaymentUseCase(repo interfaces_repo.PaymentRepository, order_Repo interfaces_repo.OrderRepository, cfg config.Config) interfaces.PaymentUseCase {
	return &paymentUseCase{
		paymentRepository: repo,
		order_Repo:        order_Repo,
		cfg:               cfg,
	}

}

func (pu *paymentUseCase) PaymentMethodID(order_id int) (int, error) {
	id, err := pu.paymentRepository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (pu *paymentUseCase) AddPaymentMethod(payment models.NewPaymentMethod) (models.PaymentDetails, error) {
	exists, err := pu.paymentRepository.CheckIfPaymentMethodAlreadyExists(payment.PaymentName)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	if exists {
		return models.PaymentDetails{}, errors.New("payment method already exists")
	}
	paymentadd, err := pu.paymentRepository.AddPaymentMethod(payment)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	return paymentadd, nil
}

// Razorpay
func (pu *paymentUseCase) MakePaymentRazorpay(orderId, userId int) ([]models.CombinedOrderDetails, string, error) {
	if orderId <= 0 || userId <= 0 {
		return nil, "", errors.New("Please provide valid IDs")
	}

	order, err := pu.order_Repo.GetOrder(orderId)
	if err != nil {
		err = errors.New("error in getting order details through order id: " + err.Error())
		return nil, "", err
	}

	client := razorpay.NewClient(pu.cfg.RazorPay_key_id, pu.cfg.RazorPay_key_secret)

	data := map[string]interface{}{
		"amount":   int(order.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return nil, "", err
	}

	razorPayOrderId, ok := body["id"].(string)
	if !ok {
		return nil, "", errors.New("failed to extract Razorpay order ID")
	}

	err = pu.paymentRepository.AddRazorPayDetails(orderId, razorPayOrderId)
	if err != nil {
		return nil, "", err
	}

	body2, err := pu.order_Repo.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, "", err
	}

	return []models.CombinedOrderDetails{body2}, razorPayOrderId, nil
}

func (pu *paymentUseCase) SavePaymentDetails(paymentId, razorId, orderId string) error {
	staus, err := pu.paymentRepository.GetPaymentStatus(orderId)
	if err != nil {
		return err
	}
	if !staus {
		err = pu.paymentRepository.UpdatePaymentDetails(razorId, paymentId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("alredy paid")

}
