package helper

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidationMessage struct {
	IsValid bool
	Message string
}

func ValidateEmail(p string) ValidationMessage {
	validate = validator.New()
	err := validate.Var(p, "required,email")

	if err != nil {
		return ValidationMessage{IsValid: false, Message: "Invalid Email"}
	}
	return ValidationMessage{IsValid: true, Message: ""}
}

func ValidatePassword(s string) ValidationMessage {
	var number, upper, special, letter bool

	length := len(s)
	if length < 8 || length > 12 {
		return ValidationMessage{IsValid: false, Message: "Password is not strong.  Max length 12 and mininum 8"}
	}

	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letter = true
		}
	}

	if number == false || upper == false || special == false || letter == false {
		return ValidationMessage{IsValid: false, Message: "Password is not strong. Must include upper and lowercase, special character, number."}
	}

	return ValidationMessage{IsValid: true}

}

func ValidatePhoneNumber(p string) ValidationMessage {
	validate = validator.New()
	err := validate.Var("+90"+p, "required,e164")

	if err != nil {
		return ValidationMessage{IsValid: false, Message: "Invalid Phone Number"}
	}
	return ValidationMessage{IsValid: true}
}

func ValidateDate(p string) ValidationMessage {
	validate = validator.New()
	err := validate.Var(p, "required,datetime=2006-01-02")

	if err != nil {
		return ValidationMessage{IsValid: false, Message: "Invalid Date format"}
	}
	return ValidationMessage{IsValid: true}
}
func ValidateAge(p time.Time) ValidationMessage {
	age := CalculateAge(p)

	if age < 18 {
		return ValidationMessage{IsValid: false, Message: "Invalid Age"}
	}
	return ValidationMessage{IsValid: true}
}
func ValidateStruct(p interface{}) ValidationMessage {
	validate = validator.New()
	err := validate.Struct(p)
	if err != nil {
		return ValidationMessage{IsValid: false, Message: string(err.Error())}
	}
	return ValidationMessage{IsValid: true}
}

func ValidateTC(tcnumber string) ValidationMessage {

	runes := []rune(tcnumber)

	if isSame(runes) {
		return ValidationMessage{IsValid: false, Message: "Invalid TC No"}
	}

	if len(runes) != 11 {
		return ValidationMessage{IsValid: false, Message: "Invalid TC No"}

	}

	odd, even, sum, rebuild := 0, 0, 0, ""

	for i := 0; i < len(runes)-2; i++ {

		a, _ := strconv.Atoi(string(runes[i]))

		if string(runes[0]) == "0" {
			return ValidationMessage{IsValid: false, Message: "Invalid TC No"}

		}

		if (i+1)%2 == 0 {
			odd += a
		} else {
			even += a
		}

		rebuild += string(runes[i])

		sum += a
	}

	ten := (even*7 - odd) % 10

	indexTen, _ := strconv.Atoi(string(runes[9]))

	eleven := (sum + indexTen) % 10

	build := string(rebuild) + strconv.Itoa(ten) + strconv.Itoa(eleven)

	if build == tcnumber {
		return ValidationMessage{IsValid: true}
	}

	return ValidationMessage{IsValid: false, Message: "Invalid TC No"}

}

func isSame(a []rune) bool {
	b := a[0:10]
	for i := 1; i < len(b); i++ {
		if b[i] != b[0] {
			return false
		}
	}
	return true
}

func ValidateCardNumber(p string) ValidationMessage {
	validate = validator.New()
	err := validate.Var(p, "required,credit_card")

	if err != nil {
		return ValidationMessage{IsValid: false, Message: "Invalid Card Number"}
	}
	return ValidationMessage{IsValid: true}
}

func ValidateCardExpDate(m int, y int) ValidationMessage {
	if isExpired(m, y) {
		return ValidationMessage{IsValid: false, Message: "Invalid Expiration Date"}
	}
	return ValidationMessage{IsValid: true}
}

func validExpiryMonth(p int) bool {
	if p < 1 || 12 < p {
		return false
	}
	return true
}

func validExpiryYear(p int) bool {
	if p < 1900 || p > 2200 {
		return false
	}
	return true
}

func isExpired(m int, y int) bool {
	if !validExpiryMonth(m) || !validExpiryYear(y) {
		return true
	}

	date := fmt.Sprintf("%d-%d-01", y, m)
	parsetime, _ := time.Parse("2006-1-02", date)

	return parsetime.Before(time.Now())
}
func ValidateLoadMoney(p float64) ValidationMessage {
	validate = validator.New()
	if p <= 0 || p > 5000 {
		return ValidationMessage{IsValid: false, Message: "Money cannot be greater than 5.000 or zero"}
	}

	err := validate.Var(p, "required,numeric")

	if err != nil {
		return ValidationMessage{IsValid: false, Message: "Invalid Money format"}
	}
	return ValidationMessage{IsValid: true, Message: ""}
}
func IsEmptyStruct(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		} else {
			return false
		}
	}
	return false
}
