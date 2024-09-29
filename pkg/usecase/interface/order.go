package interfaces

import "ShowTimes/pkg/utils/models"

type OrderUseCase interface {
	Checkout(userID int) (models.CheckoutDetails, error)
}
