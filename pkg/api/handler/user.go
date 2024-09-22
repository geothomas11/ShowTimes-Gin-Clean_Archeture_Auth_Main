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

//Google Authentication

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

	// fmt.Println("userdata", userData)

	// userDetails := models.UserDetailsGoogleAuth{}

	var userdata models.UserDetails

	json.Unmarshal(userData, &userdata)
	fmt.Println("userdata", userdata)

	userResponse, err := h.userUseCase.UserSignUp(userdata)
	if err != nil {
		c.JSON(400, err)
	}
	c.JSON(http.StatusOK, userResponse)
}

//Profile Updates

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
