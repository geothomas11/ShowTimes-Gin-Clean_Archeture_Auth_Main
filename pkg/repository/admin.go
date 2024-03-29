package repository

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}
func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Users, error) {
	var adminCompareDetails domain.Users
	if err := ad.DB.Raw("select * from users where email=?", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Users{}, err
	}
	return adminCompareDetails, nil

}
func (ad *adminRepository) GetUserByID(id int) (domain.Users, error) {
	var users domain.Users
	if err := ad.DB.Raw("select * from users where id=?", id).Scan(&users).Error; err != nil {
		return domain.Users{}, err
	}
	return users, nil

}

//	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)

func (ad *adminRepository) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {
	// pagination purpose -
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * 2
	var Get_Users []models.UserDetailsAtAdmin
	if err := ad.DB.Raw("select id,name,email,phone,blocked from users limit ? offset ?", 2, offset).Scan(&Get_Users).Error; err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return Get_Users, nil

}
func (ad *adminRepository) UpdateBlockUserByID(user models.UpdateBlock) error {

	err := ad.DB.Exec("update users set blocked = ? where id = ?", user.Blocked, user.ID).Error
	if err != nil {
		return err
	}
	return nil

}
func (ad *adminRepository) IsUserExist(id int) (bool, error) {

	var count int
	err := ad.DB.Raw("select count(*) from users where id = ?", id).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil

}
