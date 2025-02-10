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

// SendOTP sends an OTP to a provided phone number.
// @Summary Send OTP
// @Description Sends an OTP (One-Time Password) to the provided phone number for verification.
// @Tags OTP
// @Accept json
// @Produce json
// @Param OTPdata body models.OTPdata true "Phone number to send OTP"
// @Success 200 {object} response.Response  "Success: OTP sent successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in wrong format or OTP not sent"
// @Router /user/otplogin [post]
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

// VerifyOTP verifies the provided OTP code.
// @Summary Verify OTP
// @Description Verifies the provided OTP (One-Time Password) code for user authentication.
// @Tags OTP
// @Accept json
// @Produce json
// @Param VerifyData body models.VerifyData true "Data containing OTP for verification"
// @Success 200 {object} response.Response  "Success: OTP verified successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided are in wrong format or could not verify OTP"
// @Router /user/verifyotp [post]
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
