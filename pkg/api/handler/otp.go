package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"ShowTimes/pkg/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpUseCase interfaces.OtpUseCase
}

func NewOtpHandler(usecase interfaces.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		otpUseCase: usecase,
	}
}

func (ot *OtpHandler) SendOTP(c *gin.Context) {
	var phone models.OTPdata
	if err := c.BindJSON(&phone); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := ot.otpUseCase.SendOTP(phone.PhoneNumber)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "OTP not sent", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "OTP send successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}
func (ot *OtpHandler) VerifyOTP(c *gin.Context) {

	var code models.VerifyData
	if err := c.BindJSON(&code); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Couldnot verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println("at verify OTP")
	users, err := ot.otpUseCase.VerifyOTP(code)
	fmt.Println("after verifyOTP usecase")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "OTP verified successfully", users, nil)
	c.JSON(http.StatusOK, successRes)

}
