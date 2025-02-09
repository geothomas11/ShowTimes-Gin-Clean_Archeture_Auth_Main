package usecase

import (
	walletRep "ShowTimes/pkg/repository/interfaces"
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/models"
)

type WalletUsecase struct {
	walletRepo walletRep.WalletRepository
}

func NewWalletUsecase(walletRep walletRep.WalletRepository) interfaces.WalletUsecase {
	return &WalletUsecase{walletRepo: walletRep}
}
func (wu *WalletUsecase) GetWallet(userID int) (models.WalletAmount, error) {
	return wu.walletRepo.GetWallet(userID)

}
