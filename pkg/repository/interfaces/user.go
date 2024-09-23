package interfaces

import "ShowTimes/pkg/utils/models"

type UserRepository interface {
	UserSignup(user models.UserDetails) (models.UserDetailsResponse, error)
	CheckUserAvialiablity(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email string) (bool, error)
	AddAddress(userID int, address models.AddressInfoResponse) (models.AddressInfoResponse, error)
	CheckIfUserAddress(userID int) bool
	CheckUserById(userID int) bool
	ShowUserDetails(userID int) (models.UsersProfileDetails, error)
	GetAllAddress(userID int) ([]models.AddressInfoResponse, error)
	EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error)
	ChangePassword(userID, Password string) error
	FindUserById(id string) (models.UserSignInResponse, error)
}
