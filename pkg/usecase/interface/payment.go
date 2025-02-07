package interfaces

import "ShowTimes/pkg/utils/models"

type PaymentUseCase interface {
	PaymentMethodID(orderID int) (int, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error)
	MakePaymentRazorpay(orderId, userId int) (models.CombinedOrderDetails, string, error)
	SavePaymentDetails(paymentId, razorId, orderId string) error
}
