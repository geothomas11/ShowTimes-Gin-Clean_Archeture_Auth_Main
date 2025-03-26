package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"fmt"
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

// Checkout processes the checkout for the user's order.
// @Summary Process checkout
// @Description Processes the checkout for the user's order.
// @Tags Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response "Success: Checkout completed successfully"
// @Failure 400 {object} response.Response "Bad request: Getting user ID failed"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Checkout failed"
// @Router /user/orders/checkout [post]post]
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

// OrderItems places an order with items from the user's cart.
// @Summary Place order from cart
// @Description Places an order with items from the user's cart.
// @Tags Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param OrderFromCart body models.OrderFromCart true "Order details from cart"
// @Success 200 {object} response.Response "Success: Order placed successfully"
// @Failure 400 {object} response.Response "Bad request: Error in getting user ID or invalid request"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not place the order"
// @Router /user/orders [post]
func (oh *OrderHandler) OrderItems(c *gin.Context) {
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

// GetOrderDetails retrieves order details for a user.
// @Summary Retrieve order details
// @Description Retrieves order details for a user based on the provided pagination parameters.
// @Tags Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id header string true "User ID"
// @Param page query integer false "Page number (default: 1)"
// @Param count query integer false "Number of items per page (default: 10)"
// @Success 200 {object} response.Response  "Success: Retrieved order details successfully"
// @Failure 400 {object} response.Response  "Bad request: Page number or count not in correct format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve order details"
// @Router /user/orders [get]
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

// CancelOrder cancels an order by ID from an admin perspective.
// @Summary Cancel an order as admin
// @Description Cancels an order based on the provided order ID from an admin perspective.
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param order_id path integer true "Order ID to cancel"
// @Success 200 {object} response.Response "Success: Order canceled successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid order ID"
// @Failure 500 {object} response.Response "Internal server error: Could not cancel the order"
// @Router /admin/orders/{order_id} [delete]
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

// GetAllOrderDetailsForAdmin retrieves all order details for admin with pagination.
// @Summary Retrieve all order details for admin
// @Description Retrieves all order details for admin with pagination based on the provided parameters.
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param size query integer false "Number of items per page (default: 10)"
// @Success 200 {object} response.Response  "Success: Retrieved all order details for admin successfully"
// @Failure 400 {object} response.Response  "Bad request: Page number or count not in correct format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve order details for admin"
// @Router /admin/orders [get]
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

// ApproveOrder approves an order by its ID.
// @Summary Approve order
// @Description Approves an order based on the provided order ID.
// @Tags Order
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param order_id query integer true "Order ID to approve"
// @Success 200 {object} response.Response  "Success: Order approved successfully"
// @Failure 400 {object} response.Response  "Bad request: Error from orderID or couldn't approve the order"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Failed to approve the order"
// @Router /admin/orders [patch]
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

// CancelOrderFromAdmin cancels an order by its ID from an admin perspective.
// @Summary Cancel an order as an admin
// @Description Cancels an order based on the provided order ID from an admin perspective.
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param order_id query integer true "Order ID to cancel"
// @Success 200 {object} response.Response "Success: Order canceled successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid order ID"
// @Failure 500 {object} response.Response "Internal server error: Could not cancel the order"
// @Router /admin/orders [delete]
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

// ReturnOrder initiates the return process for a specific order.
// @Summary Initiate order return
// @Description Initiates the return process for an order based on the provided order ID and user ID.
// @Tags Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param order_id query integer true "Order ID to initiate return"
// @Success 200 {object} response.Response "Success: Order returned successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid order ID"
// @Failure 500 {object} response.Response "Internal server error: Couldn't process the order return"
// @Router /user/orders [patch]
func (oh *OrderHandler) ReturnOrder(c *gin.Context) {
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
		return
	}
	userID := userId.(int)
	err = oh.orderUseCase.ReturnOrder(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "coudn't cancel the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order returned Successfully", nil, nil)
	c.JSON(http.StatusOK, success)

}

// PrintInvoice generates and returns a PDF invoice for a given order.
// @Summary Generate Invoice PDF
// @Description Generates and returns a PDF invoice for the specified order ID.
// @Tags Orders
// @Accept json
// @Produce application/pdf
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param order_id query int true "Order ID"
// @Success 200 {file} application/pdf "Success: Invoice PDF file generated successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid order ID or missing parameter"
// @Failure 502 {object} response.Response "Bad gateway: Error generating or processing the invoice"
// @Router /orders/invoice [get]
func (O *OrderHandler) PrintInvoice(c *gin.Context) {
	userId, _ := c.Get("user_id")
	userID := userId.(int)

	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		err = errors.New("error in coverting order id" + err.Error())
		errRes := response.ClientResponse(http.StatusBadGateway, "error in reading the order id", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	pdf, err := O.orderUseCase.PrintInvoice(orderIdInt, userID)
	fmt.Println("error ", err)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing the invoice", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment;filename=invoice.pdf")

	pdfFilePath := "salesReport/invoice.pdf"

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing invoice", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
	c.Header("Content-Type", "application/pdf")

	c.File(pdfFilePath)

	err = pdf.Output(c.Writer)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "error in printing invoice", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "the request was succesful", pdf, nil)
	c.JSON(http.StatusOK, successRes)
}
