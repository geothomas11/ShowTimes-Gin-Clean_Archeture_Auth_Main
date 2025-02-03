package usecase

import (
	"ShowTimes/pkg/domain"
	interfaces_helper "ShowTimes/pkg/helper/interface"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"errors"
	"fmt"
	"strconv"

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
		fmt.Println("1")
		return domain.TokenAdmin{}, err

	}
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))

	if err != nil {
		fmt.Println("2")
		return domain.TokenAdmin{}, err
	}
	var AdminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&AdminDetailsResponse, &adminCompareDetails)
	if err != nil {
		fmt.Println("3")
		return domain.TokenAdmin{}, err
	}
	access, _, err := ad.helper.GenerateTokenAdmin(AdminDetailsResponse)
	if err != nil {
		fmt.Println("4")
		return domain.TokenAdmin{}, err
	}
	return domain.TokenAdmin{
		Admin:       AdminDetailsResponse,
		AccessToken: access,
		// RefreshToken: refresh,
	}, nil

}
func (ad *adminUseCase) BlockUser(id string) error {
	ID, _ := strconv.Atoi(id)
	userExist, err := ad.adminRepository.IsUserExist(ID)
	if err != nil {
		return err
	}
	if !userExist {
		return errors.New("user not exist")
	}

	user, err := ad.adminRepository.GetUserByID(ID)
	if err != nil {
		return err
	}
	if user.IsAdmin {
		return errors.New("admin's id cannot be blocked")
	}
	fmt.Println("id:", ID)
	fmt.Println("user:", user)
	var user_Blocked models.UpdateBlock

	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user_Blocked.Blocked = true
	}
	user_Blocked.ID = int(user.ID)

	err = ad.adminRepository.UpdateBlockUserByID(user_Blocked)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) UnBlockUser(id string) error {
	ID, _ := strconv.Atoi(id)
	user, err := ad.adminRepository.GetUserByID(ID)
	if err != nil {
		return err
	}
	var user_Unblock models.UpdateBlock
	if user.Blocked {
		user_Unblock.Blocked = false
	} else {
		return errors.New("already unblocked")
	}
	user_Unblock.ID = int(user.ID)
	err = ad.adminRepository.UpdateBlockUserByID(user_Unblock)
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

//admin Dashboard

func (au *adminUseCase) AdminDashboard() (models.CompleteAdminDashboard, error) {
	userDetails, err := au.adminRepository.DashboardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := au.adminRepository.DashboardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}

	orderDetails, err := au.adminRepository.DashboardOrderDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := au.adminRepository.DashboardAmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevenueDetails, err := au.adminRepository.DashboardTotalRevenueDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardAmount:  amountDetails,
		DashboardRevenue: totalRevenueDetails,
	}, nil

}

//sales Report

func (ah *adminUseCase) FilteredSalesReport(timePeriod string) (models.SalesReport, error) {

	startTime, endTime := ah.helper.GetTimeFromPeriod(timePeriod)
	salesReport, err := ah.adminRepository.FilteredSalesReport(startTime, endTime)

	if err != nil {
		return models.SalesReport{}, err
	}
	return salesReport, nil
}
