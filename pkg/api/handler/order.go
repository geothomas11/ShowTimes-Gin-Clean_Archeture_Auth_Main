package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase   interfaces.OrderUseCase
	paymentUseCase interfaces.PaymentUseCase
}

func NewOrderHandler(OUsecase interfaces.OrderUseCase, PUseCase interfaces.PaymentUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase:   OUsecase,
		paymentUseCase: PUseCase,
	}
}

func (oh *OrderHandler) Checkout(c *gin.Context) {

	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Getting ID Failed", nil, errors.New("failed to get id").Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	checkOutResp, err := oh.orderUseCase.Checkout(userID.(int))

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Checkout failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, "Successfully completed", checkOutResp, nil)
	c.JSON(http.StatusOK, successResp)

}

func (oh *OrderHandler) OrderItemsFromCart(c *gin.Context) {
	id, errs := c.Get("id")
	if !errs {
		err := errors.New("error getting id")
		errorResp := response.ClientResponse(http.StatusBadRequest, "bad request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}
	userID := id.(int)
	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorResp := response.ClientResponse(http.StatusBadRequest, "Payment option is not cod", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}
	OrderSuccessResponse, err := oh.orderUseCase.OrderItemsFromCart(orderFromCart, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully created the order", OrderSuccessResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

func (oh *OrderHandler) PlaceOrderCODD(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errr := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errr)
		return
	}

	paymentMethodID, err := oh.paymentUseCase.PaymentMethodID(order_id)
	if err != nil {
		err := response.ClientResponse(http.StatusInternalServerError, "Error from payment id", nil, err.Error())
		c.JSON(http.StatusInternalServerError, err)
		return

	}
	if paymentMethodID == 1 {
		err := oh.orderUseCase.ExecutePurchaseCOD(order_id)
		if err != nil {
			errorResp := response.ClientResponse(http.StatusInternalServerError, "error in cash  on delivery", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorResp)
			return
		}
		successRes := response.ClientResponse(http.StatusOK, "Order placed on cash on delivery", nil, nil)
		c.JSON(http.StatusOK, successRes)
	}
	if paymentMethodID != 1 {
		successRes := response.ClientResponse(http.StatusOK, "Cannot place order payment in not COD", nil, nil)
		c.JSON(http.StatusOK, successRes)

	}

}
