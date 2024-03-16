package usecase

import (
	"ShowTimes/pkg/domain"
	interfaces_helper "ShowTimes/pkg/helper/interface"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"errors"

	"ShowTimes/pkg/utils/models"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository interfaces_repo.AdminRepository
	helper          interfaces_helper.Helper
}

func NewAdminUseCase(repo interfaces_repo.AdminRepository, h interfaces_helper.Helper) interfaces.AdminUseCase {

	return &adminUseCase{
		adminRepository: repo,
		helper:          h,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {
	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))

	if err != nil {
		return domain.TokenAdmin{}, err
	}
	var AdminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&AdminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	access, refresh, err := ad.helper.GenerateTokenAdmin(AdminDetailsResponse)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	return domain.TokenAdmin{
		Admin:        AdminDetailsResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}
func (ad *adminUseCase) BlockUser(id string) error {

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) UnBlockUser(id string) error {

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}
func (ad *adminUseCase) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {

	userDetails, err := ad.adminRepository.GetUsers(page)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}
