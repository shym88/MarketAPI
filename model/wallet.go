package model

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"marketapi/database"
	"time"
)

const (
	Deposit  int = 1
	Withdraw     = 2
)

// Wallet represents data about wallet
type Wallet struct {
	ID            int `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime `gorm:"index"`
	UserID        int          `gorm:"index:idxw_userid,index,unique,not null"`
	User          User         `gorm:"references:ID"`
	CurrentAmount float64
	Status        bool `gorm:"default:true"`
	UniqueKey     uuid.UUID
}

// WalletTransaction represents data about wallet request
type WalletTransaction struct {
	ID              int `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime `gorm:"index"`
	UserID          int          `gorm:"index:idxwt_userid,index,not null"`
	User            User         `gorm:"references:ID"`
	WalletID        int          `gorm:"index:idxwt_walletid,index,not null"`
	Wallet          Wallet       `gorm:"references:ID"`
	CardID          int          `gorm:"index:idxwt_cardid,index,not null"`
	Card            Card         `gorm:"references:ID"`
	TransactionType int          `gorm:"not null"`
	Amount          float64      `gorm:"not null"`
}

// WalletTransactionReqBody represents data about wallet request
type WalletTransactionReqBody struct {
	WalletID        int     `json:"wallet_id,omitempty" binding:"required" validate:"required"`
	CardNumber      string  `json:"card_number,omitempty" binding:"required" validate:"required"`
	CVC             int     `json:"cvc,omitempty" binding:"required" validate:"required"`
	TransactionType int     `json:"transaction_type,omitempty" binding:"required" validate:"required"`
	Amount          float64 `json:"amount,omitempty" binding:"required" validate:"required"`
}

// WalletQRReqBody represents data about wallet request
type WalletQRReqBody struct {
	UserID    int       `json:"user_id,omitempty" binding:"required" validate:"required"`
	WalletID  int       `json:"wallet_id,omitempty" binding:"required" validate:"required"`
	UniqueKey uuid.UUID `json:"unique_key,omitempty" binding:"required" validate:"required"`
}

// QRCodeRes represents data about wallet request
type QRCodeRes struct {
	Code string `json:"code,omitempty"`
}

// WalletRes represents data about wallet response
type WalletRes struct {
	ID            int       `json:"id,omitempty"`
	UserID        int       `json:"user_id,omitempty"`
	CurrentAmount float64   `json:"current_amount,omitempty"`
	UniqueKey     uuid.UUID `json:"unique_key,omitempty"`
}

func (wallet *Wallet) Save() (*Wallet, error) {
	err := database.Database.Create(&wallet).Error
	if err != nil {
		return &Wallet{}, err
	}
	return wallet, nil

}
func (walletTransaction *WalletTransaction) Save(transactionType int) (*WalletTransaction, error) {
	var err error
	if transactionType == 1 {
		err = database.Database.Create(&walletTransaction).Error
	} else if transactionType == 2 {
		err = database.Database.Select("user_id", "wallet_id", "transaction_type", "amount").Create(&walletTransaction).Error
	} else {
		return &WalletTransaction{}, err
	}
	fmt.Println("BBBBBBBBBBBBBBB")
	fmt.Println(err)
	if err != nil {
		return &WalletTransaction{}, err
	}
	if !updateCurrentAmount(walletTransaction.WalletID, walletTransaction.Amount, walletTransaction.TransactionType, walletTransaction.UserID) {
		fmt.Println("cccccc")
		var walletTransactionDelete WalletTransaction
		database.Database.Where("id = ?", walletTransaction.ID).Find(&walletTransactionDelete)
		database.Database.Delete(&walletTransactionDelete)

		return &WalletTransaction{}, err
	}

	return walletTransaction, nil

}
func updateCurrentAmount(walletID int, amount float64, transactionType int, user_id int) bool {
	var wallet Wallet

	if walletID == 0 || user_id == 0 {
		return false
	}

	if err := database.Database.Where("ID = ? AND user_id = ?", walletID, user_id).Find(&wallet).Error; err != nil {
		return false
	}
	if (Wallet{}) == wallet {
		return false
	}
	var current_amount float64
	if transactionType == 1 {
		current_amount = wallet.CurrentAmount + amount
	} else if transactionType == 2 {
		current_amount = wallet.CurrentAmount - amount
	} else {
		return false
	}
	uid := uuid.New()
	result := database.Database.Model(&wallet).Where("id = ?", walletID).Select("current_amount", "unique_key").Updates(Wallet{CurrentAmount: current_amount, UniqueKey: uid})
	//result := database.Database.Model(&wallet).Update("current_amount", current_amount)
	if result.Error != nil {
		return false
	}
	return true
}
func FindWalletByUserID(user_id int) (Wallet, error) {

	var wallet Wallet
	//initializers.DB.Find(&post, id)
	err := database.Database.Where("user_id = ? AND status=true", user_id).Find(&wallet).Error
	if err != nil || (Wallet{}) == wallet {
		return Wallet{}, err
	}
	return wallet, nil

}

func FindWalletByUserIDUniqueID(walletID int, user_id int, uniqueKey uuid.UUID) (Wallet, error) {

	var wallet Wallet
	//initializers.DB.Find(&post, id)

	err := database.Database.Where("id = ? AND user_id = ? AND unique_key = ? AND status=true", walletID, user_id, uniqueKey).Find(&wallet).Error
	if err != nil || (Wallet{}) == wallet {
		return Wallet{}, err
	}

	return wallet, nil

}
