package usecase

import (
	"ShowTimes/pkg/domain"
	repo_interface "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf/v2"
)

type orderUseCase struct {
	orderRepository   repo_interface.OrderRepository
	cartRepository    repo_interface.CartRepository
	userRepository    repo_interface.UserRepository
	paymentRepository repo_interface.PaymentRepository
	walletRepo        repo_interface.WalletRepository
	couponRepo        repo_interface.CouponRepository
}

func NewOrderUseCase(orderRepo repo_interface.OrderRepository,
	walletRepo repo_interface.WalletRepository, cartRepo repo_interface.CartRepository, userRepo repo_interface.UserRepository, paymentRepo repo_interface.PaymentRepository, couponRepo repo_interface.CouponRepository) interfaces.OrderUseCase {
	return &orderUseCase{
		orderRepository:   orderRepo,
		cartRepository:    cartRepo,
		userRepository:    userRepo,
		paymentRepository: paymentRepo,
		walletRepo:        walletRepo,
		couponRepo:        couponRepo,
	}

}

func (ou *orderUseCase) Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := ou.userRepository.GetAllAddress(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := ou.orderRepository.GetAllPaymentOption()
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	cartItems, err := ou.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := ou.cartRepository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err

	}
	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil

}

func (ou *orderUseCase) OrderItems(orderFromCart models.OrderFromCart, userID int) (models.OrderSuccessResponse, error) {

	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := ou.cartRepository.CheckCart(userID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return models.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := ou.userRepository.AddressExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if !addressExist {
		return models.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	paymentExist, err := ou.paymentRepository.PaymentExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if orderBody.CouponID > 0 {
		couponExist, err := ou.couponRepo.IsCouponExistByID(orderBody.CouponID)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
		if !couponExist {
			return models.OrderSuccessResponse{}, errors.New(errmsg.ErrCouponExistFalse)
		}

	}
	if !paymentExist {
		return models.OrderSuccessResponse{}, errors.New("payment method doesnot exist")
	}

	cartItems, err := ou.cartRepository.DisplayCart(orderBody.UserID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	total, err := ou.cartRepository.TotalAmountInCart(orderBody.UserID)
	 
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	order_id, err := ou.orderRepository.OrderItems(orderBody, total)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if err := ou.orderRepository.AddOrderProducts(order_id, cartItems); err != nil {
		return models.OrderSuccessResponse{}, err
	}

	//here placeing order
	err = ou.orderRepository.UpdateOrder(order_id)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := ou.cartRepository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}

	}
	orderSuccessResponse, err := ou.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil

}


func (or *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {
	fullOrderDetails, err := or.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

// func (ou *orderUseCase) GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error) {
// 	fullOrderDetails, err := ou.orderRepository.GetOrderDetails(userID, page, count)
// 	if err != nil {
// 		return []models.FullOrderDetails{}, err

// 	}
// 	return fullOrderDetails, nil

// }
func (ou *orderUseCase) CancelOrders(orderId int, userId int) error {
	userTest, err := ou.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is done by its user")
	}

	orderProductDetails, err := ou.orderRepository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}

	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	paymentStatus, err := ou.orderRepository.GetPaymentStatus(orderId)
	if err != nil {
		return err
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		return fmt.Errorf("this order is in %s, so no point in cancelling", shipmentStatus)
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, you can return it")
	}
	if shipmentStatus == "Delivered" {
		return errors.New("the order is delivered, you can return it")
	}
	if paymentStatus == "paid" || paymentStatus == "PAID" {
		amount, err := ou.orderRepository.GetFinalPriceOrder(orderId)
		if err != nil {
			return err
		}
		err = ou.walletRepo.AddToWallet(userId, amount)
		if err != nil {
			return err
		}
	}

	err = ou.orderRepository.CancelOrders(orderId)
	if err != nil {
		return err
	}

	// Update product quantity after cancellation
	err = ou.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}

	return nil
}

func (ou *orderUseCase) GetAllOrdersAdmin(page models.Page) ([]models.CombinedOrderDetails, error) {
	if page.Page == 0 {
		page.Page = 1
	}
	offset := (page.Page - 1) * page.Size

	orderDetail, err := ou.orderRepository.GetAllOrdersAdmin(offset, page.Size)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil
}

func (ou *orderUseCase) ApproveOrder(orderId int) error {
	ShipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}

	if ShipmentStatus == "cancelled" {
		return errors.New("the order is cancelled,cannot approve it")
	}
	if ShipmentStatus == "pending" {
		return errors.New("the order is pending, cannot approve it")
	}
	if ShipmentStatus == "delivered" {
		return errors.New("this item is already delivered")
	}

	if ShipmentStatus == "processing" {
		fmt.Println("usc order")
		err := ou.orderRepository.ApproveOrder(orderId)
		if err != nil {
			return err
		}
		return nil
	}
	if ShipmentStatus == "shipped" {
		err := ou.orderRepository.ApproveCodPaid(orderId)
		if err != nil {
			return err
		}
		return nil
	}
	if ShipmentStatus == "returned" {
		err := ou.orderRepository.ApproveCodReturn(orderId)
		if err != nil {
			return err
		}
	}
	fmt.Println("last ao")
	return nil

}

func (ou *orderUseCase) CancelOrderFromAdmin(orderId int) error {
	if orderId <= 0 {
		return errors.New("invalid order id")
	}
	ok, err := ou.orderRepository.CheckOrderID(orderId)

	if !ok {
		return errors.New("order does not exist")
	}
	if err != nil {
		return err
	}
	orderProduct, err := ou.orderRepository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}

	ShipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if ShipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled")
	}
	if ShipmentStatus == "deliverd" {
		return errors.New("the order is delivered cannot be cancelled")
	}
	err = ou.orderRepository.CancelOrders(orderId)
	if err != nil {
		return err
	}
	err = ou.orderRepository.UpdateStockOfProduct(orderProduct)
	if err != nil {
		return err
	}
	return nil
}

func (ou *orderUseCase) ReturnOrder(orderId, userId int) error {
	if orderId < 0 {
		return errors.New("invalid order id")
	}
	userTest, err := ou.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("this order is not done by the  user")
	}

	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	paymentType, err := ou.orderRepository.GetPaymentType(orderId)
	if err != nil {
		return err
	}
	if shipmentStatus == "cancelled" {
		return errors.New("the order is cancelled, cannot return it")
	}
	if shipmentStatus == "pending" {
		return errors.New("the order is pending, cannot return it")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is processing cannot return it")
	}
	if shipmentStatus == "returned" {
		return errors.New("the order is returned,cannot return it")
	}
	if shipmentStatus == "shipped" {
		return errors.New("the order is shipped ,cannot return it")
	}
	amount, err := ou.orderRepository.GetFinalPriceOrder(orderId)
	if err != nil {
		return err
	}
	if paymentType == 1 {
		if shipmentStatus == "delivered" {
			err = ou.orderRepository.ReturnOrderCod(orderId)
			if err != nil {
				return err
			}
			err = ou.walletRepo.AddToWallet(userId, amount)
			if err != nil {
				return err
			}
		}

	}
	//cod
	if paymentType == 2 {
		if shipmentStatus == "delivered" {
			err = ou.orderRepository.ReturnOrderRazorPay(orderId)
			if err != nil {
				return err
			}
			err = ou.walletRepo.AddToWallet(userId, amount)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (or *orderUseCase) PrintInvoice(orderId int) (*gofpdf.Fpdf, error) {

	if orderId < 1 {
		return nil, errors.New("enter a valid order id")
	}

	order, err := or.orderRepository.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, err
	}

	fmt.Println("order usecase ", order)

	items, err := or.orderRepository.GetItemsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	fmt.Println("items usecase", items)

	if order.ShipmentStatus != "delivered" {
		return nil, errors.New("wait for the invoice until the product is received")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 30)
	pdf.SetTextColor(31, 73, 125)
	pdf.Cell(0, 20, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.Cell(0, 10, "Customer Details")
	pdf.Ln(10)
	customerDetails := []string{
		"Name: " + order.Name,
		"House Name: " + order.HouseName,
		"Street: " + order.Street,
		"State: " + order.State,
		"City: " + order.City,
	}
	for _, detail := range customerDetails {
		pdf.Cell(0, 10, detail)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Final Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(255, 255, 255)
	for _, item := range items {
		pdf.CellFormat(40, 10, item.ProductName, "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Price, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(item.Quantity), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Price*float64(item.Quantity), 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Total Price:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(totalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	offerApplied := totalPrice - order.FinalPrice

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Offer Applied:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(offerApplied, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Final Amount:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(order.FinalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(0, 10, "Generated by Watch Hive India Pvt Ltd. - "+time.Now().Format("2006-01-02 15:04:05"))
	pdf.Ln(10)

	return pdf, nil
}
