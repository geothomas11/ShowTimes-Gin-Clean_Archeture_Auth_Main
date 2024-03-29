package models

type AdminLogin struct {
	Email    string `json:"email,omitempty" validateP:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}

type AdminDetailsResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type UpdateBlock struct {
	ID      int
	Blocked bool
}
