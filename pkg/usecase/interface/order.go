package interfaces

import (
	"ShowTimes/pkg/utils/models"

	"github.com/jung-kurt/gofpdf/v2"
)

type OrderUseCase interface {
	Checkout(userID int) (models.CheckoutDetails, error)
	OrderItems(orderFromCart models.OrderFromCart, userID int) (models.OrderSuccessResponse, error)
	// ExecutePurchaseCOD(orderID int) error
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	CancelOrders(orderID int, userId int) error
	GetAllOrdersAdmin(page models.Page) ([]models.CombinedOrderDetails, error)
	ApproveOrder(orderId int) error
	CancelOrderFromAdmin(orderId int) error
	ReturnOrder(orderId, userId int) error
	// GetPaymentType(orderID int) (int, error)
	// GetPaymentStatus(orderID int) (string, error)
	// GetFinalPriceOrder(orderID int) (float64, error)
	PrintInvoice(orderId int) (*gofpdf.Fpdf, error)
}
