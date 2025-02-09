package interfaces

import "ShowTimes/pkg/utils/models"

type WalletUsecase interface {
	GetWallet(userID int) (models.WalletAmount, error)
}
