package handler

import (
	// "ShowTimes/pkg/helper/interfaces"
	"ShowTimes/pkg/config"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase interfaces.UserUseCase
}
type Response struct {
	ID      uint   `copier:"must"`
	Name    string `copier:"must"`
	Surname string `copier:"must"`
}

func NewUserHandler(usecase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}

}

// UserSignUp registers a new user.
//
// @Summary Register a new user
// @Description Registers a new user with the provided details.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param request body models.UserDetails true "User details for sign-up"
// @Success 201 {object} response.Response "User signed up successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input format or constraints not satisfied"
// @Router /user/signup [post]
func (u *UserHandler) UserSignUp(c *gin.Context) {
	var user models.UserDetails

	//Bind the user details to struct

	if err := c.BindJSON(&user); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	//check weater the data send by user gas all the correct constarins specified by Users struct

	err := validator.New().Struct(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	//USING BUSINESS LOGINC INSIDE THIS FUNCTION
	userCreated, err := u.userUseCase.UserSignUp(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "user could not signed up", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, "User signed Up successfully", userCreated, nil)
	c.JSON(http.StatusOK, successResp)
}

// LoginHandler handles user login.
//
// @Summary User login
// @Description Authenticates a user using provided credentials.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param request body models.UserLogin true "User login credentials"
// @Success 200 {object} response.Response "User logged in successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid credentials or unable to log in"
// @Router /user/login [post]
func (u *UserHandler) LoginHandler(c *gin.Context) {
	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	err := validator.New().Struct(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	user_details, err := u.userUseCase.LoginHandler(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "user could not be logged in", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "user signed successfully", user_details, nil)
	c.JSON(http.StatusOK, successResp)

}

// Authv2 initiates Google authentication.
//
// @Summary Google login
// @Description Redirects the user to Google login page.
// @Tags Users
// @Success 302 "Redirects to Google login page"
// @Router /user/google_login [get]
func (h *UserHandler) Authv2(c *gin.Context) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusSeeOther, url)
}

func (h *UserHandler) GoogleCallback(c *gin.Context) {
	fmt.Println("call back work well ")
	state := c.Query("state")
	if state != "randomstate" {
		c.JSON(http.StatusBadRequest, "state is not matching")
		return
	}

	code := c.Query("code")

	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, "make code excahnge lead to error")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, "get some access token follow to error")
		return
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "fetch user data followt ot error")

		return
	}

	var userdata models.UserDetails

	json.Unmarshal(userData, &userdata)
	fmt.Println("userdata", userdata)

	userResponse, err := h.userUseCase.UserSignUp(userdata)
	if err != nil {
		c.JSON(400, err)
	}
	c.JSON(http.StatusOK, userResponse)
}

// AddAddress adds a new address for the user.
//
// @Summary Add user address
// @Description Adds an address for the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param request body models.AddressInfoResponse true "User address details"
// @Success 200 {object} response.Response "Address added successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input or unable to add address"
// @Router /user/profile/address [post]
func (u *UserHandler) AddAddress(c *gin.Context) {
	var address models.AddressInfoResponse

	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)
	fmt.Println("id,.,.,..,.,..", userId)

	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := c.BindJSON(&address); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	err := validator.New().Struct(address)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	adrRep, err := u.userUseCase.AddAddress(userId, address)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "can not add address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, "address added  successfully", adrRep, nil)
	c.JSON(http.StatusOK, successResp)

}

// ShowUserDetails retrieves user details.
//
// @Summary Get user details
// @Description Retrieves the details of the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id header int true "User ID of the authenticated user"
// @Success 200 {object} response.Response "Successfully retrieved user details"
// @Failure 400 {object} response.Response "Bad request: Invalid user ID format or unable to retrieve details"
// @Router /user/details [get]
func (u *UserHandler) ShowUserDetails(c *gin.Context) {
	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)
	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	userResp, err := u.userUseCase.ShowUserDetails(userId)
	if err != nil {
		errREsp := response.ClientResponse(http.StatusBadRequest, "Cannot get details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "successfully got details", userResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// GetAllAddress retrieves all addresses of the authenticated user.
//
// @Summary Get all user addresses
// @Description Retrieves all saved addresses of the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id header int true "Authenticated User ID"
// @Success 200 {object} response.Response "Successfully retrieved all user addresses"
// @Failure 400 {object} response.Response "Bad request: Invalid user ID format or unable to retrieve addresses"
// @Router /user/profile/alladdress [get]

func (u *UserHandler) GetAllAddress(c *gin.Context) {
	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)
	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	userRep, err := u.userUseCase.GetAllAddress(userId)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "cannot get Addresses", nil, err)
		c.JSON(http.StatusBadRequest, errResp)
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully Got All Addresses", userRep, nil)
	c.JSON(http.StatusOK, successResp)
}

// EditProfile updates the profile details of the authenticated user.
//
// @Summary Update user profile
// @Description Updates the profile information of the authenticated user.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id header int true "Authenticated User ID"
// @Param request body models.UsersProfileDetails true "Updated user profile details"
// @Success 200 {object} response.Response "User profile updated successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input format or unable to update profile"
// @Router /user/profile [patch]
func (u *UserHandler) EditProfile(c *gin.Context) {
	var details models.UsersProfileDetails

	userString, er := c.Get("id")
	if !er {
		errResp := response.ClientResponse(http.StatusBadRequest, "Failed to get user id", nil, er)
		c.JSON(http.StatusBadRequest, errResp)
		return

	}
	userid, ers := userString.(int)
	if !ers {
		errResp := response.ClientResponse(http.StatusBadRequest, "conversion error", nil, ers)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := c.BindJSON(&details); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provide in weong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	details.ID = uint(userid)

	userResp, err := u.userUseCase.EditProfile(details)
	if err != nil {
		fmt.Println("errore", err)
		errResp := response.ClientResponse(http.StatusBadRequest, "cannot update profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully updated profile", userResp, nil)
	c.JSON(http.StatusOK, successResp)

}

// ChangePassword updates the password of the authenticated user.
//
// @Summary Change user password
// @Description Allows the authenticated user to change their password.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id header int true "Authenticated User ID"
// @Param request body models.ChangePassword true "New password details"
// @Success 200 {object} response.Response "Password changed successfully"
// @Failure 400 {object} response.Response "Bad request: Invalid input format or unable to change password"
// @Router /user/password [patch]
func (u *UserHandler) ChangePassword(c *gin.Context) {
	var change models.ChangePassword
	userString, er := c.Get("id")
	if !er {
		erresp := response.ClientResponse(http.StatusBadRequest, "Failed yo get user id", nil, er)
		c.JSON(http.StatusBadRequest, erresp)
		return
	}
	userid, ers := userString.(int)
	if !ers {
		errResp := response.ClientResponse(http.StatusBadRequest, "conversion error", nil, ers)
		c.JSON(http.StatusBadRequest, errResp)

	}
	if err := c.BindJSON(&change); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provide in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	change.UserID = uint(userid)

	err := u.userUseCase.ChangePassword(change)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "cannot change password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Password changed successfully", nil, nil)
	c.JSON(http.StatusOK, successResp)

}
