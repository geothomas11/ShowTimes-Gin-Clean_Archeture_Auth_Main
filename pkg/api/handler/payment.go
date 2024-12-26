package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase interfaces.PaymentUseCase
}

func NewPaymentHandler(usecase interfaces.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: usecase,
	}
}

func (ph *PaymentHandler) AddPaymentMethod(c *gin.Context) {
	var payment models.NewPaymentMethod

	err := c.BindJSON(&payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot add the payment method payment name", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	paymentResp, err := ph.paymentUseCase.AddPaymentMethod(payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "cannot add payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully added", paymentResp, nil)
	c.JSON(http.StatusOK, successResp)

}
