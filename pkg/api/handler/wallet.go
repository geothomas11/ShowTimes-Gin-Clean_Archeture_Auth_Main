package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletUsecase interfaces.WalletUsecase
}

func NewWalletHandler(usecase interfaces.WalletUsecase) *WalletHandler {
	return &WalletHandler{WalletUsecase: usecase}

}
func (wh *WalletHandler) GetWallet(c *gin.Context) {
	userId, exist := c.Get("id")
	if !exist {
		errs := response.ClientResponse(http.StatusBadRequest, "userId not found", nil, "user_id required")
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, ok := userId.(int)
	if !ok {
		errs := response.ClientResponse(http.StatusBadRequest, "invaliad user_id type", nil, "user_id must be an integer")
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	walletDetails, err := wh.WalletUsecase.GetWallet(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrive details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Wallet details", walletDetails, nil)
	c.JSON(http.StatusOK, success)
}
