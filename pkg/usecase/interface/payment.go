package interfaces

import "ShowTimes/pkg/utils/models"

type PaymentUseCase interface {
	PaymentMethodID(order_id int) (int, error)
	AddPaymentMethod(payment models.NewPaymentMethod) (models.PaymentDetails, error)
}
