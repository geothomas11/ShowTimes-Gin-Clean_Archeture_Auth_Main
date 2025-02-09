package usecase

import (
	walletRep "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
	"errors"
)

type WalletUsecase struct {
	walletRepo walletRep.WalletRepository
}

func NewWalletUsecase(walletRep walletRep.WalletRepository) interfaces.WalletUsecase {
	return &WalletUsecase{walletRepo: walletRep}
}
func (wu *WalletUsecase) GetWallet(userID int) (models.WalletAmount, error) {
	ok, err := wu.walletRepo.IsWalletExist(userID)
	if err != nil {
		return models.WalletAmount{}, errors.New("error in database")
	}
	if !ok {
		err = wu.walletRepo.CreateWallet(userID)
		if err != nil {
			return models.WalletAmount{}, err
		}
	}
	return wu.walletRepo.GetWallet(userID)

}
