package helper

import (
	"errors"
	"fmt"
	"marketapi/model"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type MyJWTClaims struct {
	*jwt.RegisteredClaims
}

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))
var privateKeyRefresh = []byte(os.Getenv("JWT_REFRESH_PRIVATE_KEY"))
var accesstoken_exp = getEnvAsInt("JWT_EXPIRE", 1800)
var refreshtoken_exp = getEnvAsInt("JWT_EXPIRE_REFRESH", 86400)

func GenerateJWT(user model.User) (model.TokenRes, error) {
	var err error
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims = &MyJWTClaims{&jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1000000000 * time.Duration(accesstoken_exp))),
		Subject:   strconv.Itoa(user.ID),
	},
	}

	jwt := model.TokenRes{}

	jwt.AccessToken, err = token.SignedString(privateKey)

	if err != nil {
		return jwt, err
	}
	jwtRefresh, _ := generateRefreshToken(user, jwt)
	return jwtRefresh, nil
}
func GenerateJWTForRefresh(user model.User) (model.AccessTokenRes, error) {
	var err error
	var accessTokenRes model.AccessTokenRes

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = &MyJWTClaims{&jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1000000000 * time.Duration(refreshtoken_exp))),
		Subject:   strconv.Itoa(user.ID),
	},
	}

	accessTokenRes.AccessToken, err = token.SignedString(privateKey)
	if err != nil {
		return model.AccessTokenRes{}, err
	}

	return accessTokenRes, nil
}

func generateRefreshToken(user model.User, currentToken model.TokenRes) (model.TokenRes, error) {
	var err error
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = &MyJWTClaims{&jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Subject:   strconv.Itoa(user.ID),
	},
	}
	//return token.SignedString(privateKey)

	currentToken.RefreshToken, err = token.SignedString(privateKeyRefresh)

	if err != nil {
		return currentToken, err
	}
	return currentToken, nil

}
func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid token provided")
}
func ValidateRefreshToken(tokenString string) (string, error) {
	token, err := getTokenRefresh(tokenString)
	if err != nil {
		return "", err
	}

	var userId string
	claims, ok := token.Claims.(jwt.MapClaims)
	if val, e := claims["sub"]; e {
		userId = val.(string)

	}

	//_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return userId, nil
	}

	return "", errors.New("invalid token provided")
}
func CurrentUser(context *gin.Context) (model.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return model.User{}, err
	}

	var userId string

	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	if val, ok := claims["sub"]; ok {
		userId = val.(string)

	}

	i, _ := strconv.Atoi(userId)
	user, err := model.FindUserById(i)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}
func getTokenRefresh(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
