package interfaces

import (
	"ShowTimes/pkg/utils/models"
	"mime/multipart"
	"time"

	"github.com/xuri/excelize/v2"
)

type Helper interface {
	GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error)
	GenerateTokenClients(user models.UserDetailsResponse) (string, error)
	PasswordHashing(string) (string, error)
	CompareHashAndPassword(a string, b string) error
	Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse, error)

	TwilioSetup(username string, password string)
	TwilioSendOTP(phone string, serviceID string) (string, error)
	TwilioVerifyOTP(serviceID string, code string, phone string) error

	AddImageToAwsS3(file *multipart.FileHeader) (string, error)

	ValidatePhoneNumber(phone string) bool
	ValidatePin(pin string) bool
	GetTimeFromPeriod(timePeriod string) (time.Time, time.Time)
	ValidateAlphabets(data string) (bool, error)
	ValidateDate(dateString string) bool


	ConvertToExel(sales []models.OrderDetailsAdmin, filename string) (*excelize.File, error)
}
