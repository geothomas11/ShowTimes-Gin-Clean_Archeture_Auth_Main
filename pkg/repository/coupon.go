package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type couponRepo struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &couponRepo{DB: db}
}

func (cr *couponRepo) AddCoupon(coupon models.Coupon) (models.CouponResp, error) {
	var CouponResp models.CouponResp
	dateNow := time.Now()
	query := `insert into coupons
	 (coupon_name, offer_percetage, expire_data, created_at) values (?,?,?,?) returning * `

	result := cr.DB.Raw(query, coupon.CouponName, coupon.ExpireDate, dateNow).Scan(&CouponResp)
	if result.Error != nil {
		return models.CouponResp{}, errors.New(errmsg.ErrWriteDB)
	}
	return CouponResp, nil

}
