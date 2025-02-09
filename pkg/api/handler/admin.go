package handler

import (
	"ShowTimes/pkg/helper"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	adminUseCase interfaces.AdminUseCase
}

func NewAdminHandler(usecase interfaces.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}

}

// LoginHandler handles the login operation for an admin.
// @Summary Admin login
// @Description Authenticate an admin and get access token
// @Tags admin
// @Accept json
// @Produce json
// @Param body body models.AdminLogin true "Admin credentials for login"
// @Success 200 {object} response.Response "Admin login successful"
// @Failure 400 {object} response.Response "Invalid request or constraints not satisfied"
// @Failure 401 {object} response.Response "Unauthorized: cannot authenticate user"
// @Router /admin/ [post]

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
// @Description Blocks a user based on the provided ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id query string true "User ID to block" Format(uuid)
// @Success 200 {object} response.Response "User blocked successfully"
// @Failure 400 {object} response.Response "Failed to block user"
// @Router /admin/block [put]

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
// @Description Unblocks a user based on the provided ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id query string true "User ID to unblock"
// @Success 200 {object} response.Response "User unblocked successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to unblock user"
// @Router /admin/users/unblock [patch]

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

// GetUsers retrieves users based on the provided page number.
// @Summary Retrieve users with pagination
// @Description Retrieves users based on the provided page number
// @Tags admin
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

// AdminDashBoard retrieves the dashboard information for admin.
// @Summary Retrieve admin dashboard information
// @Description Retrieves dashboard information for admin
// @Tags admin
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

// salesReport

// FilteredSalesReport retrieves the sales report for a specified time period.
// @Summary Retrieve sales report for a specific time period
// @Description Retrieves sales report for the specified time period
// @Tags admin
// @Accept json
// @Produce json
// @Param period query string true "Time period for sales report"
// @Success 200 {object} response.Response "Sales report retrieved successfully"
// @Failure 500 {object} response.Response "Unable to retrieve sales report"
// @Router /admin/currentsalesreport [get]
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

// Report by date
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
