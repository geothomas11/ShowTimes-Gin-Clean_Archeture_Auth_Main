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

// google outh

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
