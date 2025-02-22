package usecase

import (
	"ShowTimes/pkg/config"
	helper_interfaces "ShowTimes/pkg/helper/interface"
	services "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
	"ShowTimes/pkg/utils/models"
	"errors"
	"strconv"

	interfaces "ShowTimes/pkg/repository/interfaces"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo   interfaces.UserRepository
	cfg        config.Config
	helper     helper_interfaces.Helper
	walletRepo interfaces.WalletRepository
}

func NewUserUseCase(repo interfaces.UserRepository, cfg config.Config, h helper_interfaces.Helper, wallet interfaces.WalletRepository) services.UserUseCase {
	return &userUseCase{
		userRepo:   repo,
		cfg:        cfg,
		helper:     h,
		walletRepo: wallet,
	}
}

var InternalError = "Internal Server Error"
var ErrorHashingPassword = "Error In hashiing Password"

func (u *userUseCase) UserSignUp(user models.UserDetails) (models.TokenUsers, error) {
	if user.Password == "" {
		return models.TokenUsers{}, errors.New("password cannot be empty")
	}
	if user.Name == "" {
		return models.TokenUsers{}, errors.New("name cannot be empty")
	}

	phoneNumber := u.helper.ValidatePhoneNumber(user.Phone)
	if !phoneNumber {
		return models.TokenUsers{}, errors.New("invalid phone number")
	}

	userExist := u.userRepo.CheckUserAvialiablity(user.Email)
	if userExist {
		return models.TokenUsers{}, errors.New("user already exist, sign in")
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New("password does not match")
	}

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New("error hashing password")
	}
	user.Password = hashedPassword

	userData, err := u.userRepo.UserSignup(user)
	if err != nil {
		return models.TokenUsers{}, err
	}

	id := uuid.New().ID()
	str := strconv.Itoa(int(id))
	userReferral := str[:8]
	err = u.userRepo.NewReferralEntry(userData.Id, userReferral)
	if err != nil {
		return models.TokenUsers{}, errors.New("referral creation failed")
	}

	err = u.walletRepo.CreateWallet(userData.Id)
	if err != nil {
		return models.TokenUsers{}, err
	}

	if user.ReferralCode != "" {
		referredId, err := u.userRepo.GetUserIdFromReferralCode(user.ReferralCode)
		if err != nil {
			return models.TokenUsers{}, errors.New(errmsg.ErrGetDB)
		}
		if referredId != 0 {
			referralAmount := 150
			err := u.userRepo.UpdateReferralAmount(float64(referralAmount), referredId, userData.Id)
			if err != nil {
				return models.TokenUsers{}, err
			}

			amount, err := u.userRepo.AmountInRefferals(userData.Id)
			if err != nil {
				return models.TokenUsers{}, err
			}

			wallectExist, err := u.walletRepo.IsWalletExist(referredId)
			if err != nil {
				return models.TokenUsers{}, err
			}
			if !wallectExist {
				err = u.walletRepo.CreateWallet(referredId)
				if err != nil {
					return models.TokenUsers{}, err
				}
			}

			err = u.walletRepo.AddToWallet(referredId, amount)
			if err != nil {
				return models.TokenUsers{}, err
			}

			err = u.walletRepo.AddToWallet(userData.Id, float64(referralAmount))
			if err != nil {
				return models.TokenUsers{}, err
			}
		}
	}

	// Generate token and return TokenUsers
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}

	return models.TokenUsers{
		Users: userData,
		Token: tokenString,
	}, nil
}

func (u *userUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {
	ok := u.userRepo.CheckUserAvialiablity(user.Email)
	if !ok {
		return models.TokenUsers{}, errors.New("the user do not exist")
	}

	isBlocked, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New("internal error")
	}
	if isBlocked {
		return models.TokenUsers{}, errors.New("user is blocked")
	}
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUsers{}, errors.New("password is not correct")
	}
	err = u.helper.CompareHashAndPassword(user_details.Password, user.Password)

	if err != nil {
		return models.TokenUsers{}, errors.New("password incorrect")
	}

	var userDetails models.UserDetailsResponse

	userDetails.Id = int(user_details.Id)
	userDetails.Name = user_details.Name
	userDetails.Email = user_details.Email
	userDetails.Phone = user_details.Phone

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}
	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil

}

func (u *userUseCase) AddAddress(userID int, address models.AddressInfoResponse) ([]models.AddressInfoResponse, error) {
	// Validate input fields
	ok, err := u.helper.ValidateAlphabets(address.Name)
	if err != nil || !ok {
		return []models.AddressInfoResponse{}, errors.New("invalid name")
	}

	if !u.helper.ValidatePhoneNumber(address.Phone) {
		return []models.AddressInfoResponse{}, errors.New("invalid phone number")
	}

	if !u.helper.ValidatePin(address.Pin) {
		return []models.AddressInfoResponse{}, errors.New("invalid pin number")
	}

	if userID <= 0 {
		return []models.AddressInfoResponse{}, errors.New("invalid user_id")
	}

	// Check if the user exists
	if !u.userRepo.CheckUserById(userID) {
		return []models.AddressInfoResponse{}, errors.New("user does not exist")
	}

	// Add the address
	_, err = u.userRepo.AddAddress(userID, address)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	addressRep, err := u.userRepo.GetAllAddress(userID)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return addressRep, nil
}

func (u *userUseCase) ShowUserDetails(userID int) (models.UsersProfileDetails, error) {

	profile, err := u.userRepo.ShowUserDetails(userID)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return profile, nil
}

func (u *userUseCase) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {
	address, err := u.userRepo.GetAllAddress(userID)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return address, nil
}

func (u *userUseCase) EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error) {
	if user.Name == "" {
		return models.UsersProfileDetails{}, errors.New("name cannot be empty")

	}
	ok, err := u.helper.ValidateAlphabets(user.Name)
	if err != nil {
		return models.UsersProfileDetails{}, errors.New("invalid name")
	}
	if !ok {
		return models.UsersProfileDetails{}, errors.New("invalid name")
	}
	phErr := u.helper.ValidatePhoneNumber(user.Phone)
	if !phErr {
		return models.UsersProfileDetails{}, errors.New("invalid phone number")
	}
	details, err := u.userRepo.EditProfile(user)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}

	return details, nil

}

func (u *userUseCase) ChangePassword(user models.ChangePassword) error {
	if user.NewPassword == "" || user.ConfirmPassword == "" {
		return errors.New("password cannot be empty")
	}
	if user.NewPassword != user.ConfirmPassword {
		return errors.New("password mismatch")
	}
	newHashed, err := u.helper.PasswordHashing(user.NewPassword)
	if err != nil {
		return errors.New("password hashing is failed")
	}
	idString := strconv.FormatUint(uint64(user.UserID), 10)

	user_details, _ := u.userRepo.FindUserById(idString)

	err = u.helper.CompareHashAndPassword(user_details.Password, user.CurrentPassword)
	if err != nil {
		return errors.New("current password is incorrect")
	}
	err = u.userRepo.ChangePassword(idString, newHashed)
	if err != nil {
		return errors.New("password cannot change")
	}
	return nil

}
