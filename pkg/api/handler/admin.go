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

func (ad *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "User could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully blocked", nil, nil)
	c.JSON(http.StatusOK, successResp)

}

func (ad *AdminHandler) UnBlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "User couldnot be Blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "User Blocked Successfully", nil, nil)
	c.JSON(http.StatusOK, successResp)

}
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
