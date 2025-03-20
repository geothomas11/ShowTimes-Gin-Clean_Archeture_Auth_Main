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
	amount, err := wu.walletRepo.GetWallet(userID)
	if err != nil {
		return models.WalletAmount{}, err
	}
	return amount, nil
}

func (wu *WalletUsecase) GetWalletHistory(userID int) ([]models.WalletHistoryResp, error) {
	wallet, err := wu.walletRepo.GetWalletData(userID)
	if err != nil {
		return []models.WalletHistoryResp{}, err
	}
	walletResp, err := wu.walletRepo.GetWalletHistory(int(wallet.ID))
	if err != nil {
		return []models.WalletHistoryResp{}, err
	}
	return walletResp, err

}
