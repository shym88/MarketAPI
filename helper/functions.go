package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"marketapi/model"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	leapYearDay = 60
)

func CalculateAge(givenTime time.Time) int {
	presentTime := time.Now()

	switch givenLocation := givenTime.Location(); givenLocation {
	case time.UTC, nil:
		presentTime = presentTime.UTC()
	default:
		presentTime = presentTime.In(givenLocation)
	}

	givenYear := givenTime.Year()
	presentYear := presentTime.Year()

	age := presentYear - givenYear

	givenYearIsLeapYear := isLeapYear(givenYear)
	presentYearIsLeapYear := isLeapYear(presentYear)

	givenYearDay := givenTime.YearDay()
	presentYearDay := presentTime.YearDay()

	if givenYearIsLeapYear && !presentYearIsLeapYear && givenYearDay >= leapYearDay {
		givenYearDay--
	} else if presentYearIsLeapYear && !givenYearIsLeapYear && presentYearDay >= leapYearDay {
		givenYearDay++
	}

	if presentYearDay < givenYearDay {
		age--
	}

	return age
}

func isLeapYear(givenYear int) bool {
	if givenYear%400 == 0 {
		return true
	} else if givenYear%100 == 0 {
		return false
	} else if givenYear%4 == 0 {
		return true
	}

	return false
}
func Pagination(c *gin.Context) model.Pagination {
	limit := 1
	page := 1
	sort := `ID DESC`
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break
		}
	}

	return model.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func hashTo32Bytes(input string) []byte {

	data := sha256.Sum256([]byte(input))
	return data[0:]

}

// cryptoText is the text to be decrypted and the keyString is the key to use for the decryption.
// The function will output the resulting plain text string with an error variable.
func decryptString(cryptoText string, keyString string) (plainTextString string, err error) {

	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

// Takes two string, plainText and keyString.
// plainText is the text that needs to be encrypted by keyString.
// The function will output the resulting crypto text and an error variable.
func encryptString(plainText string, keyString string) (cipherTextString string, err error) {

	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func encryptAES(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func getSecretKeyForEncode() string {
	return os.Getenv("SECRETKEY_FOR_ENCODE")

}
func GenerateQRNumber(wallet model.WalletQRReqBody) (model.QRCodeRes, error) {
	var qrCodeRes model.QRCodeRes
	data, err := json.Marshal(wallet)
	if err != nil {
		return model.QRCodeRes{}, err
	}

	encrypted, err := encryptString(string(data), getSecretKeyForEncode())
	if err != nil {
		return model.QRCodeRes{}, err
	}
	qrCodeRes.Code = encrypted
	return qrCodeRes, nil

}
func DecodeQRNumber(p string) (*model.WalletQRReqBody, error) {
	// encrypt base64 crypto to original value
	dec, err := decryptString(p, getSecretKeyForEncode())
	if err != nil {
		return nil, err
	}

	wallet := model.WalletQRReqBody{}
	err = json.Unmarshal([]byte(dec), &wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
