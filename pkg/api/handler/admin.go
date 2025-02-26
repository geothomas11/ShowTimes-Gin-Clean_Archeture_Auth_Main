package handler

import (
	"ShowTimes/pkg/helper"
	interfaces "ShowTimes/pkg/helper/interface"
	interfaces_u "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
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

func NewAdminHandler(usecase interfaces_u.AdminUseCase, helper interfaces.Helper) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
		helper:       helper,
	}

}

// LoginHandler handles the login operation for an admin.
//
// @Summary Admin Login
// @Description Authenticates an admin and returns an access token.
// @Tags Admin Authentication
// @Accept json
// @Produce json
// @Param body body models.AdminLogin true "Admin login credentials"
// @Success 200 {object} response.Response "Admin login successful, returns access token"
// @Failure 400 {object} response.Response "Bad request: Invalid request or incorrect format"
// @Failure 401 {object} response.Response "Unauthorized: Incorrect credentials"
// @Failure 500 {object} response.Response "Internal server error: Login failed"
// @Router /admin/login [post]
func (ad *AdminHandler) LoginHandler(c *gin.Context) {
	var adminDetails models.AdminLogin

	if err := c.BindJSON(&adminDetails); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
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

	succesResp := response.ClientResponse(http.StatusOK, errmsg.MsgLoginSucces, admin, nil)
	c.JSON(http.StatusOK, succesResp)
}

func (ad *AdminHandler) ValidateRefreshTokenAndCreateNewAccess(c *gin.Context) {
	refreshToken := c.Request.Header.Get("RefreshToken")

	// <=CHECK THE TOKEN IS VALID OR NOT =>

	_, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("refreshscret"), nil
	})
	if err != nil {
		c.AbortWithError(500, errors.New(errmsg.ErrAccessToken))
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
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgUserBlockSucces, nil, nil)
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

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, users, nil)
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
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgsuccessAdminDash, dashbord, nil)
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
		errorResp := response.ClientResponse(http.StatusInternalServerError, "Sales report could not be retrieved", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}

	// Fix message formatting
	message := fmt.Sprintf("Successfully retrieved sales report for %s", timePeriod)
	success := response.ClientResponse(http.StatusOK, message, salesReport, nil)
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
	success := response.ClientResponse(http.StatusOK, errmsg.MsgSalesRepErr, report, nil)
	c.JSON(http.StatusOK, success)

}

// PrintSalesByDate generates and downloads a sales report for a specific date.
//
// @Summary Generate and download sales report by date
// @Description Generates and downloads a sales report in either PDF or Excel format for a given day, month, and year. If no format is specified, the report is downloaded as an Excel file by default.
// @Tags Admin Sales Reports
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param year query int true "Year of the sales report (e.g., 2025)"
// @Param month query int true "Month of the sales report (1-12)"
// @Param day query int true "Day of the sales report (1-31)"
// @Param download query string false "Download format: 'pdf' or 'excel' (defaults to 'excel')"
// @Success 200 {object} response.Response "Sales report downloaded successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid date format or unable to generate the report"
// @Failure 500 {object} response.Response "Internal server error: Unable to generate the sales report"
// @Router /admin/printsales [get]
func (a *AdminHandler) PrintSalesByDate(c *gin.Context) {
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting year", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	month := c.Query("month")
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting month", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	day := c.Query("day")
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting day", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, err := a.adminUseCase.SalesByDate(dayInt, monthInt, yearInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in getting sales details", nil, err.Error())
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

		c.Header("Content-Disposition", "attachment; filename=totalsalesreport.pdf")
		c.Header("Content-Type", "application/pdf")

		if err := pdf.Output(c.Writer); err != nil {
			errRes := response.ClientResponse(http.StatusBadGateway, "error in printing sales report", nil, err)
			c.JSON(http.StatusBadRequest, errRes)
			return
		}

	} else {
		excel, err := a.helper.ConvertToExel(body, "salesReport/sales_report.xlsx")
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

		return
	}

	succesRes := response.ClientResponse(http.StatusOK, "success", body, nil)
	c.JSON(http.StatusOK, succesRes)
}
