package usecase

import (
	"ShowTimes/pkg/domain"
	interfaces_helper "ShowTimes/pkg/helper/interface"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"

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
