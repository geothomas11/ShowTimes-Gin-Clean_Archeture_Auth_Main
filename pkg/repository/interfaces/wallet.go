package interfaces

import "ShowTimes/pkg/utils/models"

type WalletRepository interface {
	CreateWallet(userID int) error
	GetWalletData(userID int) (models.Wallet, error)
	GetWallet(userID int) (models.WalletAmount, error)
	IsWalletExist(userID int) (bool, error)
	AddToWallet(userID int, Amount float64) error
	DebitFromWallet(userID int, amount float64) error
	AddToWalletHistory(wallet models.WalletHistory) error
	GetWalletHistory(walletId int) ([]models.WalletHistoryResp, error)
	GetWalletHistoryAmount(orderId int) (float64, error)
}
