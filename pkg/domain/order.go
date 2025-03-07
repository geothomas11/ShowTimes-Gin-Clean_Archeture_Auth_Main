package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          int           `json:"user_id" gorm:"not null"`
	Users           Users         `json:"-" gorm:"foreignkey:UserID"`
	AddressID       uint          `json:"address_id" gorm:"not null"`
	Address         Address       `json:"-" gorm:"foreignkey:AddressID"`
	PaymentMethodID uint          `json:"paymentmethod_id"`
	PaymentMethod   PaymentMethod `json:"-" gorm:"foreignkey:PaymentMethodID"`
	ShipmentStatus  string        `json:"shipment_status" gorm:"default:'processing'"`
	PaymentStatus   string        `json:"payment_status" gorm:"default:'not paid'"`
	FinalPrice      float64       `json:"final_price"`
	Approval        bool          `json:"approval" gorm:"default:false"`
}

type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey;not null"`
	OrderID    uint    `json:"order_id"`
	Order      Order   `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  uint    `json:"product_id"`
	Products   Product `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type OrderSuccessResponse struct {
	OrderID        uint   `json:"order_id"`
	ShipmentStatus string `json:"shipment_status"`
}
