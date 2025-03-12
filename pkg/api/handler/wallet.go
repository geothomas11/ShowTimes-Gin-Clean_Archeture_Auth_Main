package handler

import (
	interfaces "ShowTimes/pkg/usecase/interface"
	"ShowTimes/pkg/utils/errmsg"
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

// GetWallet retrieves the wallet details of the authenticated user.
//
// @Summary Get wallet details
// @Description Retrieves the wallet details of the authenticated user.
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Wallet details retrieved successfully"
// @Failure 400 {object} response.Response "Bad request: User ID not found or invalid format"
// @Failure 500 {object} response.Response "Internal server error: Failed to retrieve wallet details"
// @Router /wallet [get]
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

// GetWalletHistory retrieves the wallet transaction history for a user.
// @Summary Get wallet transaction history
// @Description Fetches the wallet transaction history for the authenticated user.
// @Tags User Wallet Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} response.Response "Success: Retrieved wallet history successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving wallet history"
// @Failure 401 {object} response.Response "Unauthorized: Missing or invalid authentication"
// @Router /user/wallet/history [get]
func (wh *WalletHandler) GetWalletHistory(c *gin.Context) {
	userId, _ := c.Get("id")
	userID, _ := userId.(int)
	walletHistory, err := wh.WalletUsecase.GetWalletHistory(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err)
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	succesResp := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, walletHistory, nil)
	c.JSON(http.StatusOK, succesResp)
}
