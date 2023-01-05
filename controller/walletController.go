package controller

import (
	"marketapi/helper"
	"marketapi/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Wallet Details
// @Description Get Wallet Details
// @Tags wallet
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Produce json
// @Success 200 {object} model.WalletRes
// @Router /api/wallet/get [get]
func GetWalletDetails(c *gin.Context) {
	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	wallet, err := model.FindWalletByUserID(currentUser.ID)
	if (model.Wallet{}) == wallet || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card cannot found"})
		return
	}
	p := model.WalletRes{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		CurrentAmount: wallet.CurrentAmount,
		UniqueKey:     wallet.UniqueKey,
	}
	c.JSON(http.StatusOK, p)

}

// @Summary Load Money to Wallet
// @Description Load Money to Wallet
// @Tags wallet
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Produce json
// @Param WalletTransactionReqBody body model.WalletTransactionReqBody true "req params"
// @Success 200
// @Router /api/wallet/load [post]
func LoadMoneyToWallet(c *gin.Context) {

	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	var body model.WalletTransactionReqBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var errValidation = helper.ValidateStruct(body)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}
	errValidation = helper.ValidateCardNumber(body.CardNumber)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidateLoadMoney(body.Amount)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	card, err := model.FindCardByCardNumber(body.CardNumber, currentUser.ID)
	if (model.Card{}) == card || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card cannot found"})
		return
	}

	if card.UserID != currentUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Card owner is not current user",
		})
		return
	}

	errValidation = helper.ValidateCardExpDate(card.ExpMonth, card.ExpYear)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	if card.CVC != body.CVC {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid CVC",
		})
		return
	}

	walletTransaction := model.WalletTransaction{
		UserID:          currentUser.ID,
		WalletID:        body.WalletID,
		CardID:          card.ID,
		TransactionType: model.Deposit,
		Amount:          body.Amount,
	}
	data, err := walletTransaction.Save(model.Deposit)

	if (model.WalletTransaction{}) == *data || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, data)
	c.JSON(http.StatusOK, gin.H{})

}

// @Summary Generate QRCode for payment
// @Description Generate Login User QRCode for payment.
// @Tags wallet
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Produce json
// @Success 200 {object} model.QRCodeRes
// @Router /api/wallet/qr [get]
func GenerateQR(c *gin.Context) {
	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	wallet, err := model.FindWalletByUserID(currentUser.ID)
	if (model.Wallet{}) == wallet || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card cannot found"})
		return
	}

	p := model.WalletQRReqBody{
		UserID:    wallet.UserID,
		WalletID:  wallet.ID,
		UniqueKey: wallet.UniqueKey,
	}
	qrcode, _ := helper.GenerateQRNumber(p)
	if (model.QRCodeRes{}) == qrcode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QRCode cannot generated"})
		return
	}
	c.JSON(http.StatusOK, qrcode)

}
