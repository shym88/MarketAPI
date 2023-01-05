package controller

import (
	"marketapi/helper"
	"marketapi/model"
	"net/http"
	"strconv"
	"time"

	"gorm.io/datatypes"

	"github.com/gin-gonic/gin"
)

// @Summary Register
// @Description User and user wallet are created
// @Tags authentication
// @Accept  json
// @Param UserReqBody body model.UserReqBody true "req params"
// @Success 200
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var body model.UserReqBody

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	t, _ := time.Parse("2006-01-02", body.Birthday)

	var errValidation = helper.ValidateStruct(body)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidateTC(body.TC)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidateEmail(body.Email)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidatePassword(body.Password)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidatePhoneNumber(body.MobilePhone)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidateDate(body.Birthday)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}

	errValidation = helper.ValidateAge(t)

	if !errValidation.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Message,
		})
		return
	}
	user := model.User{
		TC:          body.TC,
		Name:        body.Name,
		Surname:     body.Surname,
		Email:       body.Email,
		MobilePhone: body.MobilePhone,
		Birthday:    datatypes.Date(t),
		Password:    body.Password,
	}

	_, err := user.Save()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusCreated, gin.H{"user": savedUser})
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Login
// @Description Login with email and password. Pasword must include Must include upper and lowercase, special character, number. Max length 12 and mininum 8. Returen access_token and refresh_token
// @Tags authentication
// @Accept  json
// @Produce json
// @Param LoginReqBody body model.LoginReqBody true "req params"
// @Success 200 {object} model.TokenRes
// @Router /auth/login [POST]
func Login(c *gin.Context) {
	var input model.LoginReqBody

	if err := c.ShouldBindJSON(&input); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := model.FindUserByEmail(input.Email)

	if (model.User{}) == user || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jwt)
}

// @Summary Refresh Token
// @Description Generate access_token from refresh_token
// @Tags authentication
// @Accept  json
// @Produce json
// @Param TokenRes body model.TokenRes true "req params"
// @Success 200 {object} model.AccessTokenRes
// @Router /auth/refresh [POST]
func RefreshToken(c *gin.Context) {
	var body model.TokenRes

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user model.User
	params, err := helper.ValidateRefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user.ID, err = strconv.Atoi(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWTForRefresh(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jwt)

}
