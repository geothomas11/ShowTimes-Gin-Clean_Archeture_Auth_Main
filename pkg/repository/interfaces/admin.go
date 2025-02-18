package interfaces

import (
	"ShowTimes/pkg/domain"
	"ShowTimes/pkg/utils/models"
	"time"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Users, error)
	GetUserByID(id int) (domain.Users, error)
	UpdateBlockUserByID(user models.UpdateBlock) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	IsUserExist(id int) (bool, error)

	//Admin Dashboard details
	DashboardUserDetails() (models.DashBoardUser, error)
	DashboardProductDetails() (models.DashBoardProduct, error)
	DashboardOrderDetails() (models.DashBoardOrder, error)
	DashboardTotalRevenueDetails() (models.DashBoardRevenue, error)
	DashboardAmountDetails() (models.DashBoardAmount, error)

	FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error)

	SalesByDay(yearInt int, monthInt int, dayInt int) ([]models.OrderDetailsAdmin, error)
	SalesByMonth(yearInt int, monthInt int) ([]models.OrderDetailsAdmin, error)
	SalesByYear(yearInt int) ([]models.OrderDetailsAdmin, error)
}
