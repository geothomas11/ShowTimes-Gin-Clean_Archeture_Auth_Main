package usecase

import (
	"ShowTimes/pkg/domain"
	interfaces_helper "ShowTimes/pkg/helper/interface"
	interfaces_repo "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"errors"
	"fmt"
	"strconv"
	"time"

	"ShowTimes/pkg/utils/models"

	"github.com/jinzhu/copier"

	"github.com/jung-kurt/gofpdf/v2"
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

func (au *adminUseCase) FilteredSalesReport(timePeriod string) (models.SalesReport, error) {

	if timePeriod == "" {
		err := errors.New("field cannot be empty")
		return models.SalesReport{}, err
	}
	if timePeriod != "week" && timePeriod != "month" && timePeriod != "year" {
		err := errors.New("invalid time period, available options : week, month & year")
		return models.SalesReport{}, err
	}

	startTime, endTime := au.helper.GetTimeFromPeriod(timePeriod)
	salesReport, err := au.adminRepository.FilteredSalesReport(startTime, endTime)

	if err != nil {
		return models.SalesReport{}, err
	}
	return salesReport, nil
}
func (au *adminUseCase) ExecuteSalesReportByDate(startDate, endDate string) (models.SalesReport, error) {
	parsedStartDate, err := time.Parse("02-01-2006", startDate)
	if err != nil {
		err := errors.New("enter the date in correct format")
		return models.SalesReport{}, err
	}
	isValid := !parsedStartDate.IsZero()
	if !isValid {
		err := errors.New("enter date in correct format & valid date")
		return models.SalesReport{}, err
	}
	parsedEndDate, err := time.Parse("02-01-2005", endDate)
	if err != nil {
		err := errors.New("enter the date in correct format")
		return models.SalesReport{}, err
	}
	isValid = !parsedEndDate.IsZero()
	if !isValid {
		err := errors.New("enter the date in correct format & vallid date")
		return models.SalesReport{}, err
	}
	if parsedStartDate.After(parsedEndDate) {
		err := errors.New("start date is after end date")
		return models.SalesReport{}, err
	}
	orders, err := au.adminRepository.FilteredSalesReport(parsedStartDate, parsedEndDate)
	if err != nil {
		return models.SalesReport{}, errors.New("report fetching failed")
	}
	return orders, nil

}
func (ad *adminUseCase) SalesByDate(dayInt int, monthInt int, yearInt int) ([]models.OrderDetailsAdmin, error) {

	if dayInt == 0 && monthInt == 0 && yearInt == 0 {
		return []models.OrderDetailsAdmin{}, errors.New("must enter a value for day, month, and year")
	}

	if dayInt < 0 || monthInt < 0 || yearInt < 0 {
		return []models.OrderDetailsAdmin{}, errors.New("no such values are allowded")
	}

	if yearInt >= 2020 {
		if monthInt == 0 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByYear(yearInt, monthInt, dayInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase year", body)
			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByMonth(yearInt, monthInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase month", body)
			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt > 0 && dayInt <= 31 {

			body, err := ad.adminRepository.SalesByDay(yearInt, monthInt, dayInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase day", body)
			return body, nil
		}
	}

	return []models.OrderDetailsAdmin{}, errors.New("invalid date parameters")
}

func (au *adminUseCase) PrintSalesReport(sales []models.OrderDetailsAdmin) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(31, 73, 125)
	pdf.CellFormat(0, 20, "Total Sales Report", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(0, 0, 0)

	for _, item := range sales {
		pdf.CellFormat(0, 10, "Product:"+item.ProductName, "0", 1, "L", false, 0, "")
		amount := strconv.FormatFloat(item.TotalAmount, 'f', 2, 64)
		pdf.CellFormat(0, 10, "Amount Sold:$"+amount, "0", 1, "L", false, 0, "")
		pdf.Ln(5)

	}
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(150, 150, 150)

	pdf.Cell(0, 10, "Generated by Show Times India Pvt Ltd.-"+time.Now().Format("2024-02-05 16:06:02"))

	return pdf, nil

}
