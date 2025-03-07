package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	CouponUsecase interfaces.CouponUsecase
}

func NewCouponHandler(coupon interfaces.CouponUsecase) *CouponHandler {
	return &CouponHandler{CouponUsecase: coupon}

}

// AddCoupon adds a new coupon to the system.
// @Summary Add a new coupon
// @Description Allows an admin to add a new coupon by providing the necessary details.
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param coupon body models.Coupon true "Coupon details to add"
// @Success 201 {object} response.Response "Success: Coupon added successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input format"
// @Failure 401 {object} response.Response "Unauthorized: Missing or invalid authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not add the coupon"
// @Router /admin/coupons [post]
func (ch *CouponHandler) AddCoupon(c *gin.Context) {
	var coupon models.Coupon
	if err := c.BindJSON(&coupon); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	couponRep, err := ch.CouponUsecase.AddCoupon(coupon)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgCouponAddFailed, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return

	}
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgAddSuccess, couponRep, nil)
	c.JSON(http.StatusOK, successResp)

}
