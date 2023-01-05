package model

import (
	"database/sql"
	"errors"
	"marketapi/database"
	"time"
)

type CardList []*Card
type CardListDto []*CardRes

// Card represents data about card request
type Card struct {
	ID         int `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime `gorm:"index"`
	UserID     int          `gorm:"index:idxc_userid,index,not null"`
	User       User         `gorm:"references:ID"`
	CardName   string       `gorm:"not null"`
	CardNumber string       `gorm:"index:idxc_cardnumber,index,not null"`
	FullName   string       `gorm:"not null"`
	CVC        int          `gorm:"not null"`
	ExpYear    int          `gorm:"not null"`
	ExpMonth   int          `gorm:"not null"`
	BankName   string       `gorm:"not null"`
	Status     bool         `gorm:"default:true"`
}

// CardReqBody represents data about add card request
type CardReqBody struct {
	CardName   string `json:"card_name,omitempty"  binding:"required" validate:"required"`
	CardNumber string `json:"card_number,omitempty" binding:"required" validate:"required"`
	FullName   string `json:"full_name,omitempty" binding:"required" validate:"required"`
	CVC        int    `json:"cvc,omitempty" binding:"required" validate:"required"`
	ExpYear    int    `json:"exp_year,omitempty" binding:"required" validate:"required"`
	ExpMonth   int    `json:"exp_month,omitempty" binding:"required" validate:"required"`
	BankName   string `json:"bank_name,omitempty" binding:"required" validate:"required"`
}

// CardRes represents data about list cart response
type CardRes struct {
	ID         int    `json:"id,omitempty"`
	CardName   string `json:"card_name,omitempty"`
	CardNumber string `json:"card_number,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	ExpYear    int    `json:"exp_year,omitempty"`
	ExpMonth   int    `json:"exp_month,omitempty"`
}

// CardDetailRes represents data about card request
type CardDetailRes struct {
	ID         int    `json:"id,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
	CardName   string `json:"card_name,omitempty"`
	CardNumber string `json:"card_number,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	CVC        int    `json:"cvc,omitempty"`
	ExpYear    int    `json:"exp_year,omitempty"`
	ExpMonth   int    `json:"exp_month,omitempty"`
	BankName   string `json:"bank_name,omitempty"`
	Status     bool   `json:"status,omitempty"`
}

func (card *Card) Save() (*Card, error) {

	data, _ := FindCardByCardNumber(card.CardNumber, card.UserID)
	if (Card{}) != data {
		return &Card{}, errors.New("Card is predefined in user")

	}
	err := database.Database.Create(&card).Error
	if err != nil {
		return &Card{}, err
	}
	return card, nil
}

func CreateCardRes(card Card) CardRes {
	return CardRes{
		ID:         card.ID,
		CardName:   card.CardName,
		CardNumber: card.CardNumber,
		FullName:   card.FullName,
		ExpMonth:   card.ExpMonth,
		ExpYear:    card.ExpYear,
	}
}
func CreateCardResList(cards CardList) CardListDto {
	cardListDto := CardListDto{}
	for _, p := range cards {
		card := CreateCardRes(*p)
		cardListDto = append(cardListDto, &card)
	}
	return cardListDto
}
func GetAllCards(pagination Pagination, cardList []*Card, user_id int) (CardListDto, error) {
	offset := (pagination.Page - 1) * pagination.Limit

	queryBuilder := database.Database.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Where("user_id=? AND status=true", user_id)
	//queryBuilder.Model(&Card{}).Where(cardList).Find(&cardList)
	queryBuilder.Model(&Card{}).Where(cardList).Find(&cardList)

	if queryBuilder.Error != nil || len(cardList) == 0 {
		return CardListDto{}, queryBuilder.Error
	}

	cardListDto := CreateCardResList(cardList)
	return cardListDto, nil
}
func FindCardByID(id int, user_id int) (CardDetailRes, error) {

	var card CardDetailRes
	err := database.Database.Model(&Card{}).Where("id = ? AND user_id = ? AND status=true", id, user_id).Find(&card).Error
	if err != nil || (CardDetailRes{}) == card {
		return CardDetailRes{}, err
	}
	return card, nil

}
func FindCardByCardNumber(card_number string, user_id int) (Card, error) {

	var card Card
	err := database.Database.Where("card_number = ? AND user_id = ? AND status=true", card_number, user_id).Find(&card).Error
	if err != nil || (Card{}) == card {
		return Card{}, err
	}
	return card, nil

}
func UpdateCardByID(cardReq CardDetailRes, cardID int, user_id int) (*Card, error) {

	var card Card

	if err := database.Database.Where("id = ? AND user_id = ? AND status=true", cardID, user_id).Find(&card).Error; err != nil {
		return &Card{}, err
	}

	//result := database.Database.Model(&card).Where("id = ? AND status=true", cardID).Updates(cardReq)
	result := database.Database.Model(&card).Select("card_name", "card_number", "full_name", "cvc", "exp_year", "exp_month", "bank_name").Updates(Card{
		CardNumber: cardReq.CardNumber,
		CardName:   cardReq.CardName,
		FullName:   cardReq.FullName,
		CVC:        cardReq.CVC,
		ExpYear:    cardReq.ExpYear,
		ExpMonth:   cardReq.ExpMonth,
		BankName:   cardReq.BankName,
	})

	if result.Error != nil {
		return &Card{}, result.Error
	}
	return &card, nil
}
func UpdateCardStatus(cardID int, user_id int) (*Card, error) {

	var card Card

	if err := database.Database.Where("id = ? AND user_id = ? AND status=true", cardID, user_id).Find(&card).Error; err != nil {
		return &Card{}, err
	}

	result := database.Database.Model(&card).Where("id = ? AND user_id = ? AND status=true", cardID, user_id).Update("status", false)
	if result.Error != nil {
		return &Card{}, result.Error
	}
	return &card, nil
}
func DeleteCardByID(cardID int, user_id int) (*Card, error) {

	var card Card

	if err := database.Database.Where("id = ? AND user_id = ? AND status=true", cardID, user_id).Find(&card).Error; err != nil {
		return &Card{}, err
	}

	result := database.Database.Delete(&card)
	if result.Error != nil {
		return &Card{}, result.Error
	}
	return &card, nil
}
