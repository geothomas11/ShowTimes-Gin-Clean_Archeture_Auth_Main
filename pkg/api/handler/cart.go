package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"fmt"
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

func (ch *CartHandler) AddToCart(c *gin.Context) {
	var cart models.AddCart
	userID, errb := c.Get("id")
	if !errb {
		fmt.Println("err1", errb)
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, errors.New("getting user id failed"))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&cart); err != nil {
		fmt.Println("json", cart)

		errRes := response.ClientResponse(http.StatusBadRequest, "fields provide in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID, _ = userID.(int)

	cartResp, err := ch.cartUseCase.AddToCart(cart)
	if err != nil {
		fmt.Println("err4handler", err)
		erResp := response.ClientResponse(http.StatusBadRequest, "item cannot added to cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, erResp)
	
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Item successfully added to cart", cartResp, nil)
	c.JSON(http.StatusOK, successResp)

}

func (ch *CartHandler) ListCartItems(c *gin.Context) {
	userID, ers := c.Get("id")
	if !ers {
		errRsp := response.ClientResponse(http.StatusBadRequest, "cannot list products", nil, errors.New("error ins getting user id"))
		c.JSON(http.StatusBadRequest, errRsp)
		return
	}
	cartResp, err := ch.cartUseCase.ListCartItems(userID.(int))
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot list products", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully got the cart list", cartResp, nil)
	c.JSON(http.StatusOK, successResp)

}

func (ch *CartHandler) UpdateProductQuantityCart(c *gin.Context) {
	var Cart models.AddCart
	userId, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot update quantity", nil, errors.New("error getting user id"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&Cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provide in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	Cart.UserID = userId.(int)

	cartResp, err := ch.cartUseCase.UpdateProductQuantityCart(Cart)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Updation faied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully updated", cartResp, nil)
	c.JSON(http.StatusOK, successRes)

}
func (ch *CartHandler) RemoveFromCart(c *gin.Context) {
	var cart models.RemoveFromCart

	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot update quantity", nil, errors.New("error in getting user id"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID = userID.(int)
	cartResp, err := ch.cartUseCase.RemoveFromCart(cart)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Removing from cart failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully Removed", cartResp, nil)
	c.JSON(http.StatusOK, successResp)
}
