package models

type Category struct {
	Id       uint   `json:"id" gorm:"unique;not null"`
	Category string `jsom:"category" gorm:"unique;not null"`
}
type SetNewName struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
