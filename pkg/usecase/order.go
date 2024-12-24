package usecase

import (
	repo_interface "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
)

type orderUseCase struct {
	orderRepository repo_interface.OrderRepository
	cartRepository  repo_interface.CartRepository
	userRepository  repo_interface.UserRepository
}

func NewOrderUseCase(orderRepo repo_interface.OrderRepository, cartRepo repo_interface.CartRepository, userRepo repo_interface.UserRepository) interfaces.OrderUseCase {
	return &orderUseCase{
		orderRepository: orderRepo,
		cartRepository:  cartRepo,
		userRepository:  userRepo,
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

func (ou *orderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (models.OrderSuccessResponse, error) {
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
	paymentExist, err := ou.paymentRepository.paymentExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
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
	if err:=ou.orderRepository.OrderItems(order_id,cartItems);err!=nil{
		return models.OrderSuccessResponse{},err
	}
	orderSuccessResponse,err:=ou.orderRepository.GetBriefOrderDetails()

}
