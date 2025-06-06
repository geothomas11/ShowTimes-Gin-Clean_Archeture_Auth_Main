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

// AddCouponAdmin adds a new coupon to the system.
// @Summary Add a new coupon
// @Description Allows an admin to add a new coupon by providing the necessary details.
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param coupon body models.Coupon true "Coupon details to add"
// @Success 201 {object} response.Response "Success: Coupon added successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid request payload or missing fields"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication token"
// @Failure 500 {object} response.Response "Internal server error: Could not add the coupon"
// @Router /admin/coupons [post]
func (ch *CouponHandler) AddCouponAdmin(c *gin.Context) {
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

// GetCoupon Get All coupons
// @Summary Get coupons
// @Description Retrieves coupon information
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success: Retrieved coupons successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving coupons"
// @Security BearerTokenAuth
// @Router /admin/coupon [get]
func (ch *CouponHandler) GetCouponAdmin(c *gin.Context) {
	couponResp, err := ch.CouponUsecase.GetCoupon()
	if err != nil {
		ereResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, ereResp)
	}
	successRep := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, couponResp, nil)
	c.JSON(http.StatusOK, successRep)
}

// GetCoupon Get All coupons
// @Summary Get coupons
// @Description Retrieves coupon information
// @Tags User Coupon Management
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success: Retrieved coupons successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving coupons"
// @Security BearerTokenAuth
// @Router /admin/coupon [get]
func (ch *CouponHandler) GetCouponUser(c *gin.Context) {
	couponResp, err := ch.CouponUsecase.GetCoupon()
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
	}
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, couponResp, nil)
	c.JSON(http.StatusOK, successResp)

}

// EditCoupon edits a coupon.
// @Summary Edit a coupon
// @Description Edit an existing coupon
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Param coupon body models.CouponResp true "Coupon object to edit"
// @Success 200 {object} models.CouponResp "Success: Edited coupon successfully"
// @Failure 400 {object} response.Response "Bad request: Error while editing coupon"
// @Security BearerTokenAuth
// @Router /admin/coupon/ [put]
func (ch *CouponHandler) EditCoupon(c *gin.Context) {
	var coupon models.CouponResp
	if err := c.BindJSON(&coupon); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	couuponResp, err := ch.CouponUsecase.EditCoupon(coupon)
	if err != nil {
		ereResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgEditCouponFailed, nil, err.Error())
		c.JSON(http.StatusBadRequest, ereResp)
	}
	successRep := response.ClientResponse(http.StatusOK, errmsg.MsgEditCouponSuccess, couuponResp, nil)
	c.JSON(http.StatusOK, successRep)
}
