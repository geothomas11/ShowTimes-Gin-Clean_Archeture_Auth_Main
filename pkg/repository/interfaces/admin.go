package interfaces

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/utils/models"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Users, error)
	GetUserByID(id int) (domain.Users, error)
	UpdateBlockUserByID(user models.UpdateBlock) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	IsUserExist(id int) (bool, error)
}
