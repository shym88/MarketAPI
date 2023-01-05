package controller

import (
	"marketapi/helper"
	"marketapi/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Complete Payment
// @Description Complete payment with qr code and total amount
// @Tags payment
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Produce json
// @Param PaymentReq body model.PaymentReq true "req params"
// @Success 200
// @Router /api/payment/complete [post]
func CompletePayment(c *gin.Context) {
	var body model.PaymentReq

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	data, err := helper.DecodeQRNumber(body.QRCode)

	if (&model.WalletQRReqBody{}) == data || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QR Code."})
		return
	}

	wallet, err := model.FindWalletByUserIDUniqueID(data.WalletID, data.UserID, data.UniqueKey)
	if (model.Wallet{}) == wallet || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QR Code. User and walled info not matched"})
		return
	}
	if wallet.CurrentAmount < body.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient fund"})
		return
	}

	walletTransaction := model.WalletTransaction{
		UserID:          wallet.UserID,
		WalletID:        wallet.ID,
		TransactionType: model.Withdraw,
		Amount:          body.Amount,
	}

	dataTransaction, err := walletTransaction.Save(model.Withdraw)

	if (model.WalletTransaction{}) == *dataTransaction || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, data)
	c.JSON(http.StatusOK, gin.H{})

}
