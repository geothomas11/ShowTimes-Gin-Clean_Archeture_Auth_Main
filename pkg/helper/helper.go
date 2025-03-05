package helper

import (
	"ShowTimes/pkg/config"
	interfaces "ShowTimes/pkg/helper/interface"
	"ShowTimes/pkg/utils/models"
	"errors"
	"fmt"
	"unicode"

	"mime/multipart"
	"regexp"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	cfg config.Config
}

func NewHelper(config config.Config) interfaces.Helper {
	return &helper{
		cfg: config,
	}
}

type AuthCoustumClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (helper *helper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {
	accessTokenClaims := &AuthCoustumClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCoustumClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	cfg, _ := config.LoadConfig()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.Admin_AccessKey))

	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Admin_RefreshKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (h *helper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	claims := &AuthCoustumClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  "Client",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	cfg, _ := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.User_AccessKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (h *helper) PasswordHashing(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashedPassword)
	return hash, nil
}
func (h *helper) CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}

	return nil
}

func (h *helper) Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	err := copier.Copy(udr, usr)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return *udr, nil
}

//Setting up Twilio

var client *twilio.RestClient

func (h *helper) TwilioSetup(username string, password string) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})
}

func (h *helper) TwilioSendOTP(phone string, serviceID string) (string, error) {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {
		return "", err
	}
	return *resp.Sid, nil

}

func (h *helper) TwilioVerifyOTP(serviceID string, code string, phone string) error {
	to := "+91" + phone

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)

	if err != nil {

		return err
	}
	if *resp.Status == "approved" {
		return nil
	}
	return errors.New("failed to validate OTP")
}

// Added AWS S3 bucket

func (h *helper) AddImageToAwsS3(file *multipart.FileHeader) (string, error) {
	f, openErr := file.Open()
	if openErr != nil {
		fmt.Println("Error1", openErr)
		return "", openErr
	}
	defer f.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(h.cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			h.cfg.AWSAccesskeyID,
			h.cfg.AWSSecretaccesskey,
			"",
		),
	})
	if err != nil {
		fmt.Println("Error 2", err)
		return "", err
	}

	fmt.Println("config", h.cfg.AWSRegion)
	fmt.Println("config", h.cfg.AWSAccesskeyID)
	fmt.Println("config", h.cfg.AWSSecretaccesskey)

	uploader := s3manager.NewUploader(sess)
	bucketName := "myshowtimes"
	key := file.Filename // You might want to add a path or unique identifier here

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   f,
	})

	if err != nil {
		fmt.Println("Error3", err)
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
	return url, nil
}

// Vaidation of phoneNumber
func (h *helper) ValidatePhoneNumber(phone string) bool {
	phoneNumber := phone
	pattern := `^\d{10}$`
	regex := regexp.MustCompile(pattern)
	value := regex.MatchString(phoneNumber)
	return value
}

func (h *helper) ValidatePin(pin string) bool {

	match, _ := regexp.MatchString(`^\d{4}(\d{2})?$`, pin)
	return match

}

func (h *helper) GetTimeFromPeriod(timePeriod string) (time.Time, time.Time) {
	endDate := time.Now()
	if timePeriod == "week" {
		startDate := endDate.AddDate(0, 0, -1)
		return startDate, endDate
	}
	if timePeriod == "monthly" {
		startDate := endDate.AddDate(0, -1, 0)
		return startDate, endDate
	}
	if timePeriod == "year" {
		startDate := endDate.AddDate(-1, 0, 0)
		return startDate, endDate
	}
	return endDate.AddDate(0, 0, -6), endDate

}
func (h *helper) ValidateDate(dateString string) bool {

	dateLayout := "02-01-2006"

	_, err := time.Parse(dateLayout, dateString)

	return err == nil
}

func (h *helper) ValidateAlphabets(data string) (bool, error) {
	for _, char := range data {
		if !unicode.IsLetter(char) {
			return false, errors.New("data contains non-alphabetical character")
		}
	}
	return true, nil
}

func (h *helper) ConvertToExel(sales []models.OrderDetailsAdmin, filename string) (*excelize.File, error) {
	// Default filename if not provided
	if filename == "" {
		filename = "salesReport/sales_report.xlsx"
	}

	// Create a new Excel file
	file := excelize.NewFile()

	// Set headers
	file.SetCellValue("Sheet1", "A1", "Product")
	file.SetCellValue("Sheet1", "B1", "Amount Sold")

	// Create a bold style for headers
	boldStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bold style: %v", err)
	}

	// Apply bold style to headers
	file.SetCellStyle("Sheet1", "A1", "B1", boldStyle)

	// Insert sales data
	var total float64
	for i, sale := range sales {
		row := i + 2 // Start from row 2
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), sale.ProductName)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), sale.TotalAmount)
		total += sale.TotalAmount
	}

	// Add "Final Total" row
	finalRow := len(sales) + 2
	file.SetCellValue("Sheet1", fmt.Sprintf("A%d", finalRow), "Final Total")
	file.SetCellValue("Sheet1", fmt.Sprintf("B%d", finalRow), total)

	// Create a larger font style for the "Final Total" row
	largerFontStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 14, // Adjust the size as needed
			Bold: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create larger font style: %v", err)
	}

	// Apply larger font style to the "Final Total" row
	file.SetCellStyle("Sheet1", fmt.Sprintf("A%d", finalRow), fmt.Sprintf("B%d", finalRow), largerFontStyle)

	// Save the file to disk
	err = file.SaveAs(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to save Excel file: %v", err)
	}

	// Return the created file
	return file, nil
}
