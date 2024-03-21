package domain

type Category struct {
	ID       int   `json:"id" gorm:"primarykey" `
	Category string `json:"category"`
}
