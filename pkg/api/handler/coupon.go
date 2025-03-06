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

// AddCoupon handles the creation of a new coupon
// @Summary Create a new coupon
// @Description Adds a new coupon to the system
// @Tags Coupons
// @Accept json
// @Produce json
// @Param coupon body models.Coupon true "Coupon data"
// @Success 200 {object} response.Response "Coupon successfully created"
// @Failure 400 {object} response.Response "Invalid request format"
// @Failure 409 {object} response.Response "Duplicate coupon conflict"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /coupons [post]
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
