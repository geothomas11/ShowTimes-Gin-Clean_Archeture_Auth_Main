package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"net/http"
	"strconv"

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

// AddPaymentMethod adds a new payment method.
//
// @Summary Add Payment Method
// @Description Adds a new payment method using the provided details.
// @Tags Admin Payment Methods
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param request body models.NewPaymentMethod true "Details of the new payment method"
// @Success 200 {object} response.Response "Success: Payment method added successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input or payment method could not be added"
// @Router /admin/payment [post]
func (ph *PaymentHandler) AddPaymentMethod(c *gin.Context) {
	var payment models.NewPaymentMethod

	err := c.BindJSON(&payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot add the payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	paymentResp, err := ph.paymentUseCase.AddPaymentMethod(payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot add payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully added payment method", paymentResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// MakePaymentRazorpay processes a Razorpay payment.
//
// @Summary Make Payment using Razorpay
// @Description Initiates a payment process using Razorpay for the specified user and order.
// @Tags Payments
// @Accept json
// @Produce html
// @Param user_id query int true "User ID of the payer"
// @Param order_id query int true "Order ID associated with the payment"
// @Success 200 {string} string "Payment initiated successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid user ID or order ID"
// @Router /payment/razorpay [get]
func (ph *PaymentHandler) MakePaymentRazorpay(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid user ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Invalid order ID format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, razorId, err := ph.paymentUseCase.MakePaymentRazorpay(orderIdInt, userIdInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Error processing payment", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"final_price": body.FinalPrice * 100,
		"razor_id":    razorId,
		"user_id":     userId,
		"order_id":    body.OrderId,
		"user_name":   body.Name,
		"total":       int(body.FinalPrice),
	})
}

// VerifyPayment verifies and saves payment details.
//
// @Summary Verify Payment
// @Description Verifies the payment details after a successful transaction and updates the database.
// @Tags Payments
// @Accept json
// @Produce json
// @Param order_id query string true "Order ID associated with the payment"
// @Param payment_id query string true "Payment ID received from the payment gateway"
// @Param razor_id query string true "Razorpay payment identifier"
// @Success 200 {object} response.Response "Successfully updated payment details"
// @Failure 500 {object} response.Response "Internal server error: Could not update payment details"
// @Router /payment/verify [get]
func (pu *PaymentHandler) VerifyPayment(c *gin.Context) {
	orderId := c.Query("order_id")
	paymentId := c.Query("payment_id")
	razorId := c.Query("razor_id")

	if err := pu.paymentUseCase.SavePaymentDetails(paymentId, razorId, orderId); err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not update payment details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated payment details", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
