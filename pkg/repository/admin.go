package repository

import (
	"ShowTimes/pkg/domain"
	interfaces "ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"
	"time"

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
	if err := ad.DB.Raw("SELECT id,name,email,phone,blocked FROM users limit ? offset ?", 3, offset).Scan(&Get_Users).Error; err != nil {
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

// admin Dashboard
func (ar *adminRepository) DashboardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := ar.DB.Raw("select count (*) from users where is_admin ='false' ").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		err = errors.New("cannot get total users from db")
		return models.DashBoardUser{}, err
	}
	return userDetails, nil

}

func (ar *adminRepository) DashboardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := ar.DB.Raw("SELECT COUNT(*) FROM users WHERE is_admin = 'false' ").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		err = errors.New("cannot get products from db")
		return models.DashBoardProduct{}, err

	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM products WHERE stock <= 0").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		err = errors.New("cannot get stock from db")
		return models.DashBoardProduct{}, err
	}
	return productDetails, nil
}

func (ar *adminRepository) DashboardAmountDetails() (models.DashBoardAmount, error) {
	var amountDetails models.DashBoardAmount
	query := `SELECT coalesce(sum(final_price),0) FROM orders WHERE payment_status ='PAID' `
	err := ar.DB.Raw(query).Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		err = errors.New("cannot get total amount from  db")
		return models.DashBoardAmount{}, err
	}
	query = `SELECT coalesce(sum(final_price),0) FROM orders WHERE payment_status ='not_paid'
	 AND
	  shipment_status = 'pending'
	  OR 
	  shipment_status = 'processing'
	  OR 
	  shipment_status = 'shipped'
	    `
	err = ar.DB.Raw(query).Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		err = errors.New("cannot get pending amount from db")
		return models.DashBoardAmount{}, err
	}
	return amountDetails, nil
}

func (ar *adminRepository) DashboardOrderDetails() (models.DashBoardOrder, error) {
	var orderDetails models.DashBoardOrder
	err := ar.DB.Raw("SELECT count(*) FROM orders WHERE payment_status = 'PAID' ").Scan(&orderDetails.CompletedOrder).Error
	if err != nil {
		err = errors.New("cannot get total order from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM  orders WHERE shipment_status = 'pending' OR shipment_status = 'processing'").Scan(&orderDetails.PendingOrder).Error
	if err != nil {
		err = errors.New("cannoot get pending orders from db")
		return models.DashBoardOrder{}, err
	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status = 'cancelled' ").Scan(&orderDetails.CancelledOrder).Error
	if err != nil {
		err = errors.New("cannot get cancelled order from db")
		return models.DashBoardOrder{}, err

	}
	err = ar.DB.Raw("SELECT COUNT(*) FROM orders ").Scan(&orderDetails.TotalOrderItem).Error
	if err != nil {
		err = errors.New("cannot get total order items from db")
		return models.DashBoardOrder{}, err
	}
	return orderDetails, nil

}
func (ar *adminRepository) DashboardTotalRevenueDetails() (models.DashBoardRevenue, error) {
	var revenueDetails models.DashBoardRevenue
	startTime := time.Now().AddDate(0, 0, 1)
	err := ar.DB.Raw("SELECT coalesce(sum(final_price),0) FROM orders WHERE payment_status ='PAID' and created_at >=?", startTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		err = errors.New("cannot get today revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(0, -1, 1)
	err = ar.DB.Raw("SELECT COALESCE(sum(final_price),0) FROM orders WHERE payment_status = 'PAID' and created_at >= ?", startTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		err = errors.New("cannot get month revenue from db")
		return models.DashBoardRevenue{}, err
	}
	startTime = time.Now().AddDate(-1, 1, 1)
	err = ar.DB.Raw("SELECT COALESCE(sum(final_price),0) FROM orders WHERE payment_status = 'PAID' and created_at >= ?", startTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		err = errors.New("cannot get year revenue from db")
		return models.DashBoardRevenue{}, err

	}
	return models.DashBoardRevenue{}, nil
}

//Sales Report

func (ar *adminRepository) FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error) {
	var salesReport models.SalesReport

	query :=
		`SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status = 'PAID' AND created_at >= ? AND created_at <= ?`
	result := ar.DB.Raw(query, startTime, endTime).Scan(&salesReport.TotalSales)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = ar.DB.Raw("SELECT COUNT(*) FROM orders WHERE created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.TotalOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	query =
		`SELECT COUNT(*) FROM orders WHERE payment_status = 'PAID' AND created_at >= ? AND created_at <= ? `

	result = ar.DB.Raw(query, startTime, endTime).Scan(&salesReport.CompletedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	query = `SELECT COUNT(*) FROM orders WHERE shipment_status = 'processing' AND approval = 'false' AND created_at >= ? AND created_at <= ?`

	result = ar.DB.Raw(query, startTime, endTime).Scan(&salesReport.PendingOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	query = `SELECT COUNT(*) FROM orders WHERE shipment_status = 'cancelled' AND created_at >= ? AND created_at <= ?`
	result = ar.DB.Raw(query, startTime, endTime).Scan(&salesReport.CancelledOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	query = `SELECT COUNT(*) FROM orders WHERE shipment_status = 'returned' AND created_at >= ? AND created_at <= ?`
	result = ar.DB.Raw(query, startTime, endTime).Scan(&salesReport.ReturnedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error

	}

	var productID int
	query =
		`SELECT product_id FROM order_items GROUP BY product_id ORDER BY SUM(quantity) DESC LIMIT 1`

	result = ar.DB.Raw(query).Scan(&productID)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}

	result = ar.DB.Raw("SELECT product_name FROM products WHERE id = ?", productID).Scan(&salesReport.TrendingProduct)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	return salesReport, nil

}
func (ar *adminRepository) SalesByYear(yearInt int, monthInt int, dayInt int) ([]models.OrderDetailsAdmin, error) {
	var orderDetails []models.OrderDetailsAdmin

	query := `SELECT i.product_name,SUM(oi.total_price) AS total_amount FROM orders o JOIN order_items oi ON o.id = oi.order_id
              JOIN inventories i ON oi.inventory_id = i.id
              WHERE o.payment_status = 'PAID'
                AND EXTRACT(YEAR FROM o.created_at) = ?
              GROUP BY i.product_name`
	if err := ar.DB.Raw(query, yearInt).Scan(&orderDetails).Error; err != nil {
		return []models.OrderDetailsAdmin{}, err
	}
	return orderDetails, nil
}

func (ar *adminRepository) SalesByMonth(yearInt int, monthInt int) ([]models.OrderDetailsAdmin, error) {
	var orderDetails []models.OrderDetailsAdmin

	query := `SELECT i.product_name, SUM(oi.total_price) AS total_amount
              FROM orders o
              JOIN order_items oi ON o.id = oi.order_id
              JOIN inventories i ON oi.inventory_id = i.id
              WHERE o.payment_status = 'PAID'
			  AND EXTRACT(YEAR FROM o.created_at) = ?
			  AND EXTRACT(MONTH FROM o.created_at) = ?
              GROUP BY i.product_name`
	if err := ar.DB.Raw(query, yearInt, monthInt).Scan(&orderDetails).Error; err != nil {
		return []models.OrderDetailsAdmin{}, err
	}
	return orderDetails, nil

}

func (ar *adminRepository) SalesByDay(yearInt int, monthInt int, dayInt int) ([]models.OrderDetailsAdmin, error) {
	var orderDetails []models.OrderDetailsAdmin

	query := `SELECT p.product_name, SUM(oi.total_price) AS total_amount
              FROM orders o
              JOIN order_items oi ON o.id = oi.order_id
              JOIN products p ON oi.product_id = p.id
              WHERE o.payment_status = 'PAID'
			  AND EXTRACT(YEAR FROM o.created_at) = ?
			  AND EXTRACT(MONTH FROM o.created_at) = ?
                AND EXTRACT(DAY FROM o.created_at) = ?
              GROUP BY p.product_name`
	if err := ar.DB.Raw(query, yearInt, monthInt, dayInt).Scan(&orderDetails).Error; err != nil {
		return []models.OrderDetailsAdmin{}, err
	}
	return orderDetails, nil

}
