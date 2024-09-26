package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
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

func (ch *CartHandler) AddToCart(c *gin.Context) {
	var cart models.AddCart
	userID, errb := c.Get("id")
	if !errb {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, errors.New("getting user id failed"))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&cart); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provide in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID, _ = userID.(int)

	cartResp, err := ch.cartUseCase.AddToCart(cart)
	if err != nil {
		erResp := response.ClientResponse(http.StatusBadRequest, "item cannot added to cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, erResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Item successfully added to cart", cartResp, nil)
	c.JSON(http.StatusOK, successResp)

}
