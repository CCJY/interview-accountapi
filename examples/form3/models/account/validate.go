package account

import (
	"encoding/json"
	"reflect"
	"regexp"
	"time"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

type RFC3339NanoDate struct {
	time.Time
}

const DateFormat = time.RFC3339Nano

func (d RFC3339NanoDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *RFC3339NanoDate) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d RFC3339NanoDate) String() string {
	return d.Time.Format(DateFormat)
}

func IsSWiftCode(fl validator.FieldLevel) bool {
	regex := "^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$"
	SwiftCodeRegex := regexp.MustCompile(regex)

	return SwiftCodeRegex.MatchString(fl.Field().String())
}

func IsBankAccountNumber(fl validator.FieldLevel) bool {
	regex := "^(\\d){7,8}$"
	BankAccountNumberRegex := regexp.MustCompile(regex)

	return BankAccountNumberRegex.MatchString(fl.Field().String())
}

// http://ht5ifv.serprest.pt/extensions/tools/IBAN/index.html
func IsIBAN(fl validator.FieldLevel) bool {
	regex := "^GB\\d{2}[A-Z]{4}\\d{14}$"
	account, ok := fl.Parent().Interface().(AccountAttributes)
	if ok {

		switch account.Country {
		case "AU":
			return fl.Field().Interface() == nil
		case "BE":
			regex = "^BE\\d{14}$"
		case "CA":
			return fl.Field().Interface() == nil
		case "EE":
			regex = "^EE\\d{18}$"
		case "FI":
			regex = "^FI\\d{16}$"
		case "FR":
			regex = "^FR\\d{12}[0-9A-Z]{11}\\d{2}$"
		case "DE":
			regex = "^DE\\d{20}$"
		case "GR":
			regex = "^GR\\d{9}[0-9A-Z]{16}$"
		case "HK":
			return fl.Field().Interface() == nil
		case "IE":
			regex = "^IE\\d{2}[A-Z]{4}\\d{14}$"
		case "IT":
			regex = "^IT\\d{2}[A-Z]\\d{10}[0-9A-Z]{12}$"
		case "LU":
			regex = "^LU\\d{5}[0-9A-Z]{13}$"
		case "NL":
			regex = "^NL\\d{2}[A-Z]{4}\\d{10}$"
		case "PT":
			regex = "^PT\\d{23}$"
		case "ES":
			regex = "^ES\\d{22}$"
		case "SE":
			regex = "^SE\\d{22}$"
		case "CH":
			regex = "^CH\\d{7}[0-9A-Z]{12}$"
		case "US":
			return fl.Field().Interface() == nil
		case "GB":
			regex = "^GB\\d{2}[A-Z]{4}\\d{14}$"
		}
	}

	IBANRegex := regexp.MustCompile(regex)
	return IBANRegex.MatchString(fl.Field().String())
}

func IsReferenceMask(fl validator.FieldLevel) bool {
	regex := "^[#$?\\\\]{1,35}?$"
	ReferenceMaskRegex := regexp.MustCompile(regex)

	return ReferenceMaskRegex.MatchString(fl.Field().String())
}

func NewValidator() *validator.Validate {
	var validator = validator.New()
	validator.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(uuid.UUID); ok {
			return valuer.String()
		}
		return nil
	}, uuid.UUID{})

	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	// CustomerRegister(validator, trans, "swift_code", "", IsSWiftCode)
	CustomerRegister(validator, trans, "bank_account_number", "{0}{1}", IsBankAccountNumber)
	CustomerRegister(validator, trans, "iban", "{0}{1}", IsIBAN)
	// CustomerRegister(validator, trans, "reference_mask", "", IsReferenceMask)

	return validator
}

func CustomerRegister(v *validator.Validate, trans ut.Translator, tag string, text string, validFn func(f1 validator.FieldLevel) bool) {
	v.RegisterValidation(tag, validFn)
	v.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})

}
