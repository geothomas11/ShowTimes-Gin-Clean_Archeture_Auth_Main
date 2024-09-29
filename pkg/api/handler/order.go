package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
}

func NewOrderHandler(usecase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: usecase,
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
