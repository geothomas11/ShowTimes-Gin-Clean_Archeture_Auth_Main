package usecase

import (
	helper "ShowTimes/pkg/helper/interface"
	repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"errors"
)

type cuponUsecase struct {
	couponRepo repo.CouponRepository
	h          helper.Helper
}

func NewCouponUsecase(couponRep repo.CouponRepository, h helper.Helper) interfaces.CouponUsecase {
	return &cuponUsecase{
		couponRepo: couponRep,
		h:          h,
	}

}

func (cu *cuponUsecase) AddCoupon(coupon models.Coupon) (models.CouponResp, error) {
	if coupon.CouponName == "" {
		return models.CouponResp{}, errors.New(errmsg.ErrFieldEmpty)

	}
	if coupon.OfferPercentage <= 0 {
		return models.CouponResp{}, errors.New(errmsg.ErrDataZero)
	}
	formattedExpireDate := coupon.ExpireDate.Format("03-02-2006")
	ok := cu.h.ValidateDate(formattedExpireDate)
	if !ok {
		return models.CouponResp{}, errors.New(errmsg.ErrInvalidDate)
	}
	couponResp, err := cu.couponRepo.AddCoupon(coupon)
	if err != nil {
		return models.CouponResp{}, err
	}
	return couponResp, nil

}
