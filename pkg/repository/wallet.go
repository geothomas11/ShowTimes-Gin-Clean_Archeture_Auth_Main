package repository

import (
	"ShowTimes/pkg/repository/interfaces"
	"ShowTimes/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type WalletDB struct {
	Db *gorm.DB
}

func NewWalletReppsitory(db *gorm.DB) interfaces.WalletRepository {
	return &WalletDB{Db: db}
}
func (wr *WalletDB) CreateWallet(userID int) error {
	result := wr.Db.Exec("INSERT INTO wallets(userID) VALUES (?) RETURNING id  ", userID)
	if result.Error != nil {
		fmt.Println("err in wallet creation id")
		return errors.New("cannot create wallet error at database")
	}
	return nil

}

func (wr *WalletDB) GetWallet(userID int) (models.WalletAmount, error) {
	var WalletAmount models.WalletAmount
	if err := wr.Db.Raw("select amount from wallets where user_id = ?", userID).Scan(&WalletAmount).Error; err != nil {
		return models.WalletAmount{}, err
	}
	return WalletAmount, nil

}
func (wr *WalletDB) IsWalletExist(userID int) (bool, error) {
	var count int
	err := wr.Db.Raw("select coount(*) from wallets where user_id = ?", userID).Scan(&count).Error
	if err != nil {
		return false, errors.New("cannot get wallet details")
	}
	return count >= 1, nil

}
