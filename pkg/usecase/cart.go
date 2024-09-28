package usecase

import (
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"errors"
	"fmt"
)

type cartUseCase struct {
	CartRepository    interfaces_repo.CartRepository
	productRepository interfaces_repo.ProductRepository
}

func NewCartUseCase(repoc interfaces_repo.CartRepository, repop interfaces_repo.ProductRepository) interfaces.CartUseCase {
	return &cartUseCase{
		CartRepository:    repoc,
		productRepository: repop,
	}

}

func (cu *cartUseCase) AddToCart(cart models.AddCart) (models.CartResponse, error) {

	if cart.ProductID < 1 || cart.UserID < 1 {
		fmt.Println("prd", cart.ProductID)
		return models.CartResponse{}, errors.New("invalid product id or user id")
	}
	if cart.Quantity < 1 {
		fmt.Println("qunt", cart.Quantity)
		return models.CartResponse{}, errors.New("quantity must be greater")
	}
	is_avialibale, err := cu.productRepository.CheckProductAvailable(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if !is_avialibale {
		return models.CartResponse{}, errors.New("product is not available")
	}
	stock, err := cu.CartRepository.CheckStock(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if stock < int(cart.Quantity) {
		return models.CartResponse{}, err
	}
	price, err := cu.productRepository.GetPriceOfProduct(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}

	QuantityOfProductInCart, err := cu.CartRepository.QuantityOfProductInCart(cart.UserID, int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if (QuantityOfProductInCart + int(cart.Quantity)) > 20 {
		return models.CartResponse{}, errors.New("limit exceeds")
	}

	finalPrice := (price * float64(cart.Quantity))

	if QuantityOfProductInCart == 0 {
		err := cu.CartRepository.AddToCart(cart.UserID, int(cart.ProductID), int(cart.Quantity), finalPrice)
		if err != nil {
			return models.CartResponse{}, err

		}
	} else {
		currentTotal, err := cu.CartRepository.TotalPriceForProductInCart(cart.UserID, int(cart.ProductID))
		if err != nil {
			return models.CartResponse{}, err
		}
		err = cu.CartRepository.UpdateCart(QuantityOfProductInCart+int(cart.Quantity), currentTotal+finalPrice, cart.UserID, int(cart.ProductID))
		if err != nil {
			return models.CartResponse{}, err
		}

	}

	cartDetails, err := cu.CartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.CartRepository.GetTotalPrice(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}

func (cu *cartUseCase) ListCartItems(userID int) (models.CartResponse, error) {

	cartDetails, err := cu.CartRepository.DisplayCart(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.CartRepository.GetTotalPrice(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}

func (cu *cartUseCase) UpdateProductQuantityCart(cart models.AddCart) (models.CartResponse, error) {
	fmt.Println("cart P.id ", cart.ProductID)
	fmt.Println("cart quant", cart.Quantity)
	fmt.Println("cart userid", cart.UserID)
	if cart.Quantity < 1 || cart.ProductID < 1 {
		return models.CartResponse{}, errors.New("invalid product id or quantity")
	}
	is_available, err := cu.productRepository.CheckProductAvailable(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if !is_available {
		return models.CartResponse{}, errors.New("product is not available")
	}
	stock, err := cu.CartRepository.CheckStock(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}

	if stock < int(cart.Quantity) {
		return models.CartResponse{}, errors.New("out of stock")
	}

	if int(cart.Quantity) > 20 {
		return models.CartResponse{}, errors.New("limit exceeds")
	}

	err = cu.CartRepository.UpdateProductQuantityCart(cart)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartDetails, err := cu.CartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.CartRepository.GetTotalPrice(cart.UserID)
	if err != nil {

		return models.CartResponse{}, err
	}

	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil
}

func (cu *cartUseCase) RemoveFromCart(cart models.RemoveFromCart) (models.CartResponse, error) {

	if cart.ProductID < 1 {
		fmt.Println("RFC usr id", cart.UserID)
		fmt.Println("RFC usc", cart.ProductID)
		return models.CartResponse{}, errors.New("product id cannot be empty")
	}
	err := cu.CartRepository.RemoveFromCart(cart)
	if err != nil {
		return models.CartResponse{}, err
	}
	is_avialibale, err := cu.CartRepository.CheckCart(cart.UserID)
	if !is_avialibale {
		return models.CartResponse{}, err
	}
	cartDetails, err := cu.CartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.CartRepository.GetTotalPrice(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}
