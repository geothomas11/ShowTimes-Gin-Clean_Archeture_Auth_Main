package domain

import "gorm.io/gorm"

type Referral struct {
	gorm.Model
	UserID         uint    `json:"user_id" gorm:"uniquekey; not null"`
	Users          Users   `json:"-" gorm:"foreignkey:UserID"`
	ReferralCode   string  `json:"referral_code"`
	ReferralAmount float64 `json:"referral_amount"`
	ReferredUserID uint    `json:"referred_user_id"`
}
