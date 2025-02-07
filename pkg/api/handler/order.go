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

// func (oh *OrderHandler) PlaceOrderCOD(c *gin.Context) {
// 	order_id, err := strconv.Atoi(c.Query("order_id"))
// 	if err != nil {
// 		errr := response.ClientResponse(http.StatusBadRequest, "error from orderID", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errr)
// 		return
// 	}

// 	paymentMethodID, err := oh.paymentUseCase.PaymentMethodID(order_id)
// 	if err != nil {
// 		err := response.ClientResponse(http.StatusInternalServerError, "Error from payment id", nil, err.Error())
// 		c.JSON(http.StatusInternalServerError, err)
// 		return

// 	}
// 	if paymentMethodID == 1 {
// 		err := oh.orderUseCase.ExecutePurchaseCOD(order_id)
// 		if err != nil {
// 			errorResp := response.ClientResponse(http.StatusInternalServerError, "error in cash  on delivery", nil, err.Error())
// 			c.JSON(http.StatusInternalServerError, errorResp)
// 			return
// 		}
// 		successRes := response.ClientResponse(http.StatusOK, "Order placed on cash on delivery", nil, nil)
// 		c.JSON(http.StatusOK, successRes)
// 	}
// 	if paymentMethodID != 1 {
// 		successRes := response.ClientResponse(http.StatusOK, "Cannot place order payment in not COD", nil, nil)
// 		c.JSON(http.StatusOK, successRes)

// 	}

// }
func (oh *OrderHandler) GetOrderDetails(c *gin.Context) {
	pagestr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Page is not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count is not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	id, errs := c.Get("id")
	if !errs {
		err := errors.New("couldn't get id")
		erorrRes := response.ClientResponse(http.StatusInternalServerError, "Error in getting id", nil, err.Error())
		c.JSON(http.StatusInternalServerError, erorrRes)
		return

	}
	UserID := id.(int)
	OrderDetails, err := oh.orderUseCase.GetOrderDetails(UserID, page, pageSize)
	if err != nil {
		erorRes := response.ClientResponse(http.StatusInternalServerError, "could not get order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, erorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Full order details", OrderDetails, nil)
	c.JSON(http.StatusOK, successResp)

}

func (oh *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error from orderID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	id, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errRes := response.ClientResponse(http.StatusBadRequest, "error form userid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}
	userID := id.(int)
	err = oh.orderUseCase.CancelOrders(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not place order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Cancel Successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (oh *OrderHandler) GetAllOrdersAdmin(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	countStr := c.DefaultQuery("size", "10")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count is not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	var pageStruct models.Page
	pageStruct.Page = page
	pageStruct.Size = pageSize
	allOrderDetails, err := oh.orderUseCase.GetAllOrdersAdmin(pageStruct)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not retived the order details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order details retived successfully", allOrderDetails, nil)
	c.JSON(http.StatusOK, success)

}

func (oh *OrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Invalid order ID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	err = oh.orderUseCase.ApproveOrder(orderId)
	if err != nil {
		errs := response.ClientResponse(http.StatusConflict, "Order approval failed", nil, err.Error())
		c.JSON(http.StatusConflict, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Order approved successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (oh *OrderHandler) CancelOrderFromAdmin(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from order id", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = oh.orderUseCase.CancelOrderFromAdmin(order_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "could not complete the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order cancel successfully", nil, nil)
	c.JSON(http.StatusOK, success)

}
func (oh *OrderHandler) ReturnOrderCod(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not cancel orderID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userId, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errRes := response.ClientResponse(http.StatusBadRequest, "error from userid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}
	userID := userId.(int)
	err = oh.orderUseCase.ReturnOrderCod(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "coudn't cancel the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order returned Successfully", nil, nil)
	c.JSON(http.StatusOK, success)

}
