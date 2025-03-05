package interfaces

import "ShowTimes/pkg/utils/models"

type CouponUsecase interface {
	AddCoupon(coupon models.Coupon) (models.CouponResp, error)
}
