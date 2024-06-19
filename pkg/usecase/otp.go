package usecase

import (
	"ShowTimes/pkg/config"
	helper_int "ShowTimes/pkg/helper/interface"
	OTP_repo "ShowTimes/pkg/repository/interfaces"
	OTP_UseCase "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	cfg           config.Config
	otpRepository OTP_repo.OtpRepository
	helper        helper_int.Helper
}

func NewOtpUsecase(cfg config.Config, repo OTP_repo.OtpRepository, h helper_int.Helper) OTP_UseCase.OtpUseCase {
	return &otpUseCase{
		cfg:           cfg,
		otpRepository: repo,
		helper:        h,
	}

}
func (ot *otpUseCase) SendOTP(phone string) error {
	ok := ot.otpRepository.FindUserByMobileNumber(phone)

	if !ok {
		return errors.New("the user doesnot exist")
	}

	ot.helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)

	_, err := ot.helper.TwilioSendOTP(phone, ot.cfg.SERVICESID)
	if err != nil {
		return errors.New("error occured while generating OTP")
	}
	return nil
}

func (ot *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUsers, error) {

	ot.helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	err := ot.helper.TwilioVerifyOTP(ot.cfg.SERVICESID, code.Code, code.PhoneNumber)

	if err != nil {
		//this guard clause catches the error code runs only until here
		return models.TokenUsers{}, errors.New("errors while verifying")
	}

	//if user is authenticated using OTP send back user details
	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(code.PhoneNumber)
	if err != nil {
		return models.TokenUsers{}, err
	}
	tokenString, err := ot.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, err
	}
	var user models.UserDetailsResponse
	err = copier.Copy(&user, &userDetails)
	if err != nil {
		return models.TokenUsers{}, err
	}
	return models.TokenUsers{
		Users: user,
		Token: tokenString,
	}, nil
}
