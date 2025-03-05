package interfaces

import "ShowTimes/pkg/utils/models"

type CouponRepository interface {
	AddCoupon(models.Coupon) (models.CouponResp, error)
}
