package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userDatabase{DB: db}
}

func (c *userDatabase) CheckUserAvialiablity(email string) bool {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM USERS WHERE email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func (c *userDatabase) UserSignup(user models.UserDetails) (models.UserDetailsResponse, error) {

	var UserDetails models.UserDetailsResponse
	err := c.DB.Raw("insert into users (name,email,password,phone) values (?,?,?,?)returning id, name, email, phone", user.Name, user.Email, user.Password, user.Phone).Scan(&UserDetails).Error

	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return UserDetails, nil
}
func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error) {
	var user_details models.UserSignInResponse

	err := c.DB.Raw(` SELECT * FROM 
users
WHERE email = ?
and
blocked = false`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserSignInResponse{}, errors.New("error checking user details")

	}
	return user_details, nil
}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	var isBlocked bool
	err := cr.DB.Raw("select blocked from users where email = ?", email).Scan(&isBlocked).Error
	if err != nil {
		return false, err
	}
	return isBlocked, nil
}
func (db *userDatabase) CheckIfUserAddress(userID int) bool {
	var count int
	qurry := "select count(*) from addresses where user_id = $1"
	if err := db.DB.Raw(qurry, userID).Scan(&count).Error; err != nil {
		return false
	}
	return true

}
func (db *userDatabase) CheckUserById(userID int) bool {
	var count int
	qurry := "select count(*) from users where id = $1"
	if err := db.DB.Raw(qurry, userID).Scan(&count).Error; err != nil {
		return false
	}
	return true
}

func (c *userDatabase) FindUserById(id string) (models.UserSignInResponse, error) {
	var user_details models.UserSignInResponse
	err := c.DB.Raw("SELECT * FROM users  WHERE id = ?", id).Scan(&user_details).Error
	if err != nil {
		return models.UserSignInResponse{}, errors.New("error in chechking user details")
	}
	return user_details, nil

}

func (db *userDatabase) AddAddress(userID int, address models.AddressInfoResponse) (models.AddressInfoResponse, error) {

	querry := "INSERT INTO addresses(user_id,name,house_name,street,city,state,phone,pin) VALUES (?,?,?,?,?,?,?,?)"
	err := db.DB.Exec(querry, userID, address.Name, address.HouseName, address.Street, address.City, address.State, address.Phone, address.Pin).Error
	if err != nil {
		return models.AddressInfoResponse{}, errors.New("could not add address, db error")
	}
	return models.AddressInfoResponse{}, nil

}

func (db *userDatabase) ShowUserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	query := "SELECT id,name,email,phone from users where id = ?"

	result := db.DB.Raw(query, userID).Scan(&userDetails)
	if result.Error != nil {
		return models.UsersProfileDetails{}, result.Error
	}
	return userDetails, nil

}

func (db *userDatabase) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {
	qurey := "SELECT * from addresses WHERE user_id = ?"
	var address []models.AddressInfoResponse
	result := db.DB.Raw(qurey, userID).Scan(&address)
	if result.Error != nil {
		return []models.AddressInfoResponse{}, result.Error
	}
	return address, nil
}
func (db *userDatabase) EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error) {
	querry := "UPDATE users SET name=?,email=?,phone=? WHERE id=?"
	err := db.DB.Exec(querry, user.Name, user.Email, user.Phone, user.ID).Error
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return user, nil

}

func (db *userDatabase) ChangePassword(userID, password string) error {
	query := "UPDATE users SET password = ? where id = ?"
	err := db.DB.Exec(query, password, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func (or *userDatabase) AddressExist(orderBody models.OrderIncoming) (bool, error) {

	var count int
	if err := or.DB.Raw("SELECT COUNT(*) FROM addresses WHERE user_id=? AND id =?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ur *userDatabase) NewReferralEntry(userid int, userReferral string) error {
	err := ur.DB.Exec("INSERT INTO referrals (user_id, referral_code, referral_amount) VALUES (?, ?, ?)", userid, userReferral, 0).Error
	if err != nil {
		return err
	}
	return nil

}
func (ur *userDatabase) GetUserIdFromReferralCode(ReferralCode string) (int, error) {

	var referredUserId int
	err := ur.DB.Raw("SELECT user_id FROM referrals WHERE referral_code = ?", ReferralCode).Scan(&referredUserId).Error
	if err != nil {
		return 0, nil
	}

	return referredUserId, nil
}

func (ur *userDatabase) UpdateReferralAmount(referralAmount float64, referredUserId, currentUserID int) error {

	err := ur.DB.Exec("UPDATE referrals SET referral_amount = ? , referred_user_id = ? WHERE user_id = ? ", referralAmount, referredUserId, currentUserID).Error
	if err != nil {
		return err
	}

	// find the current amount in referred users referral table and add 100 with that
	err = ur.DB.Exec("UPDATE referrals SET referral_amount = referral_amount + ? WHERE user_id = ? ", referralAmount, referredUserId).Error
	if err != nil {
		return err
	}

	return nil

}
func (ur *userDatabase) AmountInRefferals(userID int) (float64, error) {
	var a float64
	err := ur.DB.Raw("SELECT referral_amount FROM referrals WHERE user_id = ?", userID).Scan(&a).Error
	if err != nil {
		return 0.0, err
	}
	return a, nil
}
