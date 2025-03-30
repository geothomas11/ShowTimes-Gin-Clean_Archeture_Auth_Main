package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase interfaces.CartUseCase
}

func NewCartHandler(usecase interfaces.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase,
	}
}

// AddToCart adds an item to the user's cart.
// @Summary Add an item to the cart
// @Description Adds an item to the user's cart based on the provided details.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param AddCart body models.AddCart true "Item details to add to the cart"
// @Success 200 {object} response.Response "Success: Item added to cart successfully"
// @Failure 400 {object} response.Response "Bad request: Fields provided in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Cannot add item to cart"
// @Router /user/cart [post]
func (ch *CartHandler) AddToCart(c *gin.Context) {
	var cart models.AddCart
	userID, errb := c.Get("id")
	if !errb {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided in the wrong format", nil, errors.New(errmsg.MsdGetIdErr))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.ErrFormat, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID, _ = userID.(int)

	cartResp, err := ch.cartUseCase.AddToCart(cart)
	if err != nil {
		erResp := response.ClientResponse(http.StatusBadRequest, "Item could not be added to cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, erResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Item successfully added to cart", cartResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// ListCartItems retrieves the list of items in the user's cart.
// @Summary Retrieve cart items
// @Description Retrieves the list of items in the user's cart based on the user ID.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response "Success: Retrieved cart items successfully"
// @Failure 400 {object} response.Response "Bad request: Cannot list products"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not retrieve the cart list"
// @Router /user/cart/list [get]
func (ch *CartHandler) ListCartItems(c *gin.Context) {
	userID, ers := c.Get("id")
	if !ers {
		errRsp := response.ClientResponse(http.StatusBadRequest, "Cannot list products", nil, errors.New("error getting user ID"))
		c.JSON(http.StatusBadRequest, errRsp)
		return
	}
	cartResp, err := ch.cartUseCase.ListCartItems(userID.(int))
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot list products", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully retrieved the cart list", cartResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// UpdateProductQuantityCart updates the quantity of a product in the user's cart.
// @Summary Update product quantity in cart
// @Description Updates the quantity of a product in the user's cart based on the provided details.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param UpdateCart body models.AddCart true "Product details to update quantity"
// @Success 200 {object} response.Response "Success: Quantity updated successfully"
// @Failure 400 {object} response.Response "Bad request: Cannot update quantity or fields are in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Update failed"
// @Router /user/cart [patch]
func (ch *CartHandler) UpdateProductQuantityCart(c *gin.Context) {
	var cart models.AddCart
	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot update quantity", nil, errors.New("error getting user ID"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID = userID.(int)

	cartResp, err := ch.cartUseCase.UpdateProductQuantityCart(cart)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Update failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully updated", cartResp, nil)
	c.JSON(http.StatusOK, successRes)
}

// RemoveFromCart removes a product from the user's cart.
// @Summary Remove product from cart
// @Description Removes a product from the user's cart based on the provided details.
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param RemoveFromCart body models.RemoveFromCart true "Product details to remove from cart"
// @Success 200 {object} response.Response "Success: Product removed from cart successfully"
// @Failure 400 {object} response.Response "Bad request: Cannot remove product or fields are in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Removing from cart failed"
// @Router /user/cart/remove [delete]
func (ch *CartHandler) RemoveFromCart(c *gin.Context) {
	var cart models.RemoveFromCart

	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot remove product", nil, errors.New("error getting user ID"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID = userID.(int)
	cartResp, err := ch.cartUseCase.RemoveFromCart(cart)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgRemoveCartErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgRemoveCartSuccess, cartResp, nil)
	c.JSON(http.StatusOK, successResp)
}
