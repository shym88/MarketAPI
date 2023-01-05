package controller

import (
	"marketapi/helper"
	"strconv"

	"github.com/gin-gonic/gin"

	"marketapi/model"
	"net/http"
)

// @Summary Add Card
// @Description Add Card for User
// @Tags cards
// @Accept  json
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Param CardReqBody body model.CardReqBody true "req params"
// @Success 200
// @Router /api/card/add [post]
func AddCard(c *gin.Context) {
	currentUser, err := helper.CurrentUser(c)
	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	var body model.CardReqBody

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

	errValidation = helper.ValidateCardExpDate(body.ExpMonth, body.ExpYear)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	card := model.Card{
		UserID:     currentUser.ID,
		CardName:   body.CardName,
		CardNumber: body.CardNumber,
		FullName:   body.FullName,
		CVC:        body.CVC,
		ExpYear:    body.ExpYear,
		ExpMonth:   body.ExpMonth,
		BankName:   body.BankName,
	}

	_, err = card.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusCreated, gin.H{"card": savedCard})
	c.JSON(http.StatusOK, gin.H{})

}

// @Summary List Card
// @Description List user's Card, pagination and sort
// @Tags cards
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Produce json
// @Success 200 {object} model.CardRes
// @Router /api/card/list [post]
func ListCard(c *gin.Context) {
	currentUser, err := helper.CurrentUser(c)
	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	var cardList []*model.Card
	//initializers.DB.Find(&cardList)
	pagination := helper.Pagination(c)

	data, err := model.GetAllCards(pagination, cardList, currentUser.ID)

	if err != nil || len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)

}

// @Summary Get Card Details
// @Description Get Card Details
// @Tags cards
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Produce json
// @Param id   path int  true  "Card ID"
// @Success 200 {object} model.CardDetailRes
// @Router /api/card/get/{id} [get]
func GetCardDetailsById(c *gin.Context) {
	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	cardId, err := strconv.Atoi(c.Param("id"))

	if cardId == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID cannot be null",
		})
		return
	}

	data, err := model.FindCardByID(cardId, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary Update Card
// @Description Update Card
// @Tags cards
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Param id   path int  true  "Card ID"
// @Param CardReqBody body model.CardReqBody true "req params"
// @Success 200
// @Router /api/update/:id [put]
func UpdateCard(c *gin.Context) {
	var body model.CardReqBody

	currentUser, err := helper.CurrentUser(c)
	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	cardId, err := strconv.Atoi(c.Param("id"))
	if cardId == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID cannot be null",
		})
		return
	}
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

	errValidation = helper.ValidateCardExpDate(body.ExpMonth, body.ExpYear)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	cardReq := model.CardDetailRes{
		CardName:   body.CardName,
		CardNumber: body.CardNumber,
		FullName:   body.FullName,
		CVC:        body.CVC,
		ExpYear:    body.ExpYear,
		ExpMonth:   body.ExpMonth,
		BankName:   body.BankName,
	}

	_, err = model.UpdateCardByID(cardReq, cardId, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}

// @Summary Deactivate Card
// @Description Deactivate Card
// @Tags cards
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Param id   path int  true  "Card ID"
// @Success 200
// @Router /api/deactivate/:id [put]
func DeactivateCard(c *gin.Context) {
	cardId, err := strconv.Atoi(c.Param("id"))

	if cardId == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID cannot be null",
		})
		return
	}
	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	_, err = model.UpdateCardStatus(cardId, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}

// @Summary Delete Card
// @Description Deactivate Card
// @Tags cards
// @Security ApiKeyAuth
// @Param Authorization header string true "write Bearer before token"
// @Accept  json
// @Param id   path int  true  "Card ID"
// @Success 200
// @Router /api/delete/:id [delete]
func DeleteCard(c *gin.Context) {

	cardId, err := strconv.Atoi(c.Param("id"))

	if cardId == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID cannot be null",
		})
		return
	}
	currentUser, err := helper.CurrentUser(c)

	if (model.User{}) == currentUser || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Session is empty",
		})
		return
	}

	_, err = model.DeleteCardByID(cardId, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
