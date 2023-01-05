package model

import (
	"database/sql"
	"html"
	"marketapi/database"
	"strings"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID          int `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime   `gorm:"index"`
	TC          string         `gorm:"index:idxu_tc,index,unique,not null"`
	Name        string         `gorm:"not null"`
	Surname     string         `gorm:"not null"`
	Email       string         `gorm:"index:idxu_email,index,unique,not null"`
	MobilePhone string         `gorm:"not null"`
	Birthday    datatypes.Date `gorm:"not null"`
	Password    string         `gorm:"not null"`
	Status      bool           `gorm:"default:true"`
}

// UserReqBody represents data about register user request
type UserReqBody struct {
	TC          string `json:"tc,omitempty" binding:"required" validate:"required"`
	Name        string `json:"name,omitempty" binding:"required" validate:"required"`
	Surname     string `json:"surname,omitempty" binding:"required" validate:"required"`
	Email       string `json:"email,omitempty" binding:"required" validate:"required"`
	MobilePhone string `json:"mobile_phone,omitempty" binding:"required" validate:"required"`
	Birthday    string `json:"birthday,omitempty" binding:"required" validate:"required"`
	Password    string `json:"password,omitempty" binding:"required" validate:"required"`
}

// LoginReqBody represents data about login request
type LoginReqBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() (*User, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}

	uid := uuid.New()
	wallet := Wallet{
		UserID:    user.ID,
		UniqueKey: uid,
	}
	_, err = wallet.Save()
	if err != nil {
		var userDelete User
		database.Database.Where("id = ?", user.ID).Find(&userDelete)
		database.Database.Delete(&userDelete)
		return &User{}, err
	}

	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return err
	}
	user.Password = string(passwordHash)
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := database.Database.Where("email=?", email).Find(&user).Error
	if err != nil || (User{}) == user {
		return User{}, err
	}
	return user, nil
}

func FindUserById(id int) (User, error) {
	var user User
	err := database.Database.Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
