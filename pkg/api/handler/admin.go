package handler

import (
	"ShowTimes/pkg/helper"
	interfaces "ShowTimes/pkg/helper/interface"
	interfaces_u "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	adminUseCase interfaces_u.AdminUseCase
	helper       interfaces.Helper
}

func NewAdminHandler(usecase interfaces_u.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}

}

// LoginHandler handles the login operation for an admin.
// @Summary Admin login
// @Description Authenticate an admin and get an access token.
// @Tags Admin Authentication
// @Accept json
// @Produce json
// @Param body body models.AdminLogin true "Admin credentials for login"
// @Success 200 {object} response.Response "Admin login successful"
// @Failure 400 {object} response.Response "Invalid request or incorrect format"
// @Failure 401 {object} response.Response "Unauthorized: cannot authenticate user"
// @Router /admin/login [post]
func (ad *AdminHandler) LoginHandler(c *gin.Context) {
	var adminDetails models.AdminLogin

	if err := c.BindJSON(&adminDetails); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "details is not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot authenticate Admin", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	c.Set("Access", admin.AccessToken)
	// c.Set("Refresh", admin.RefreshToken)

	succesResp := response.ClientResponse(http.StatusOK, "login Successfully", admin, nil)
	c.JSON(http.StatusOK, succesResp)
}

// FilteredSalesReport retrieves the sales report for a specified time period.
// @Summary Retrieve sales report for a specific time period
// @Description Retrieves the sales report based on the provided time period.
// @Tags Admin Sales Reports
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Param period query string true "Time period for sales report (e.g., daily, weekly, monthly)"
// @Success 200 {object} response.Response "Success: Sales report retrieved successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid period format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not retrieve sales report"
// @Router /admin/salesreport [get]
func (ad *AdminHandler) ValidateRefreshTokenAndCreateNewAccess(c *gin.Context) {
	refreshToken := c.Request.Header.Get("RefreshToken")

	// <=CHECK THE TOKEN IS VALID OR NOT =>

	_, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("refreshscret"), nil
	})
	if err != nil {
		c.AbortWithError(401, errors.New("refresh token is invalid:user have to login agian"))
		return
	}
	claims := &helper.AuthCoustumClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	newAccessToken, err := token.SignedString([]byte("accesssecret"))
	if err != nil {
		c.AbortWithError(500, errors.New("error in creating new access token"))
	}
	c.JSON(200, newAccessToken)
}

// BlockUser blocks a user by ID.
// @Summary Block a user
// @Description Blocks a user based on the provided ID.
// @Tags User Management
// @Accept json
// @Produce json
// @Param id query string true "User ID to block"
// @Success 200 {object} response.Response "User blocked successfully"
// @Failure 400 {object} response.Response "Failed to block user"
// @Router /admin/block [patch]
func (ad *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "User could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "User blocked Successfully ", nil, nil)
	c.JSON(http.StatusOK, successResp)

}

// UnBlockUser unblocks a user by ID.
// @Summary Unblock a user
// @Description Unblocks a user based on the provided ID.
// @Tags User Management
// @Accept json
// @Produce json
// @Param id query string true "User ID to unblock"
// @Success 200 {object} response.Response "User unblocked successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to unblock user"
// @Router /admin/unblock [patch]
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "User couldnot be Blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "User UnBlocked Successfully", nil, nil)
	c.JSON(http.StatusOK, successResp)

}

// GetUsers retrieves users with pagination.
// @Summary Retrieve users with pagination
// @Description Retrieves users based on the provided page number.
// @Tags User Management
// @Accept json
// @Produce json
// @Param page query int true "Page number for pagination"
// @Success 200 {object} response.Response "Users retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to retrieve users"
// @Router /admin/users [get]
func (ad *AdminHandler) GetUsers(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ad.adminUseCase.GetUsers(page)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)

}

// AdminDashboard retrieves dashboard information for admin.
// @Summary Retrieve admin dashboard information
// @Description Retrieves dashboard information for admin.
// @Tags Admin Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Admin dashboard retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to retrieve dashboard"
// @Router /admin/dashboard [get]
func (ah *AdminHandler) AdminDashboard(c *gin.Context) {
	dashbord, err := ah.adminUseCase.AdminDashboard()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrive records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully retrived the dashbord", dashbord, nil)
	c.JSON(http.StatusOK, successResp)

}

// FilteredSalesReport retrieves sales report for a specified time period.
// @Summary Retrieve sales report for a current time period
// @Description Retrieves sales report for the specified time period.
// @Tags Sales Reports
// @Accept json
// @Produce json
// @Param period query string true "Time period for sales report (e.g., weekly, monthly)"
// @Success 200 {object} response.Response "Sales report retrieved successfully"
// @Failure 500 {object} response.Response "Unable to retrieve sales report"
// @Router /admin/salesreport [get]
func (ah *AdminHandler) FilteredSalesReport(c *gin.Context) {
	timePeriod := c.Query("period")
	salesReport, err := ah.adminUseCase.FilteredSalesReport(timePeriod)
	if err != nil {

		errorResp := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrived", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}
	message := "current" + timePeriod + "Successfully retived sales report "
	success := response.ClientResponse(http.StatusOK, message, salesReport, err)
	c.JSON(http.StatusOK, success)
}

// SalesReportByDate retrieves sales reports based on a date range.
// @Summary Retrieve sales report by date range
// @Description Retrieves sales report based on the provided start and end dates.
// @Tags Sales Reports
// @Accept json
// @Produce json
// @Param start query string true "Start date (YYYY-MM-DD)"
// @Param end query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} response.Response "Sales report retrieved successfully"
// @Failure 500 {object} response.Response "Unable to retrieve sales report"
// @Router /admin/salesreport/daterange [get]
func (ah *AdminHandler) SalesReportByDate(c *gin.Context) {
	startDateStr := c.Query("start")
	endDateStr := c.Query("end")
	if startDateStr == "" || endDateStr == "" {
		err := response.ClientResponse(http.StatusBadGateway, "Start or end date is empty", nil, "Empty date string")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	report, err := ah.adminUseCase.ExecuteSalesReportByDate(startDateStr, endDateStr)
	if err != nil {
		errorResp := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrived", nil, err.Error)
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}
	success := response.ClientResponse(http.StatusOK, "sales report retrived successfully", report, nil)
	c.JSON(http.StatusOK, success)

}

// PrintSalesReportByDate generates and retrieves a sales report within a specified date range.
//
// @Summary Generate sales report by date range
// @Description Generates and retrieves a sales report based on the provided start and end dates.
// @Tags Admin Sales Reports
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param start query string true "Start date in YYYY-MM-DD format"
// @Param end query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} response.Response "Sales report retrieved successfully"
// @Failure 400 {object} response.Response "Bad request: Start or end date is missing"
// @Failure 500 {object} response.Response "Internal server error: Unable to retrieve sales report"
// // @Router /admin/salesreport/print [get]
// func (ah *AdminHandler) PrintSalesByDate(c *gin.Context) {
// 	startDateStr := c.Query("start")
// 	endDateStr := c.Query("end")
// 	if startDateStr == "" || endDateStr == "" {
// 		err := response.ClientResponse(http.StatusBadRequest, "start or end date is empty", nil, "Empty date string")
// 		c.JSON(http.StatusBadRequest, err)
// 		return
// 	}
// 	report, err := ah.adminUseCase.ExecuteSalesReportByDate(startDateStr, endDateStr)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusInternalServerError, "sales report could not be retrieved", nil, err.Error())
// 		c.JSON(http.StatusInternalServerError, errorRes)
// 		return
// 	}

// 	success := response.ClientResponse(http.StatusOK, "sales report retrieved successfully", report, nil)
// 	c.JSON(http.StatusOK, success)
// }

// PrintSalesByDate generates and downloads a sales report for a specific date.
//
// @Summary Generate and download sales report by date
// @Description Generates and downloads a sales report in either PDF or Excel format for the given day, month, and year.
// @Tags Admin Sales Reports
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param year query int true "Year of the sales report"
// @Param month query int true "Month of the sales report"
// @Param day query int true "Day of the sales report"
// @Param download query string false "Download format (pdf or excel), defaults to excel"
// @Success 200 {file} application/pdf "PDF sales report downloaded successfully"
// @Success 200 {file} application/vnd.openxmlformats-officedocument.spreadsheetml.sheet "Excel sales report downloaded successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid date format or unable to generate report"
// @Failure 500 {object} response.Response "Internal server error: Unable to generate sales report"
// @Router /admin/salesreport/download [get]
func (a *AdminHandler) PrintSalesByDate(c *gin.Context) {
	year := c.Query("year")
	fmt.Println("err1", year)
	yearInt, err := strconv.Atoi(year)
	fmt.Println("err2")

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting year", nil, err.Error())
		fmt.Println("err3")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	month := c.Query("month")
	fmt.Println("err4")
	monthInt, err := strconv.Atoi(month)
	fmt.Println("err5")

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting month", nil, err.Error())
		fmt.Println("err6")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	day := c.Query("day")
	fmt.Println("err7")
	dayInt, err := strconv.Atoi(day)
	fmt.Println("err8")

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting day", nil, err.Error())
		fmt.Println("err9")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, err := a.adminUseCase.SalesByDate(dayInt, monthInt, yearInt)
	fmt.Println("err10")

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting sales details", nil, err.Error())
		fmt.Println("err11")
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	download := c.Query("download")
	
	if download == "pdf" {
		pdf, err := a.adminUseCase.PrintSalesReport(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
		c.Header("Content-Disposition", "attachment;filename=totalsalesreport.pdf")

		pdfFilePath := "salesReport/totalsalesreport.pdf"

		err = pdf.OutputFileAndClose(pdfFilePath)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		c.Header("Content-Disposition", "attachment; filename=total_sales_report.pdf")
		c.Header("Content-Type", "application/pdf")

		c.File(pdfFilePath)

		c.Header("Content-Type", "application/pdf")

		err = pdf.Output(c.Writer)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	} else {
		fmt.Println("body ", body)
		excel, err := a.helper.ConvertToExel(body)
		if err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

		fileName := "sales_report.xlsx"

		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		if err := excel.Write(c.Writer); err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "Error in serving the sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}
	}

	succesRes := response.ClientResponse(http.StatusOK, "success", body, nil)
	c.JSON(http.StatusOK, succesRes)
}
