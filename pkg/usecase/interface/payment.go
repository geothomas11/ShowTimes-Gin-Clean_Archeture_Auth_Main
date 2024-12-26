package interfaces

import "ShowTimes/pkg/utils/models"

type PaymentUseCase interface {
	PaymentMethodID(orderID int) (int, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error)
}
