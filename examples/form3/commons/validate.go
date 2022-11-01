package commons

import (
	"encoding/json"
	"reflect"
	"regexp"
	"time"

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
func IsUKIBAN(fl validator.FieldLevel) bool {
	regex := "^GB\\d{2}[A-Z]{4}\\d{14}$"
	UKIBANRegex := regexp.MustCompile(regex)

	return UKIBANRegex.MatchString(fl.Field().String())
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

	validator.RegisterValidation("swift_code", IsSWiftCode)
	validator.RegisterValidation("bank_account_number", IsBankAccountNumber)
	validator.RegisterValidation("uk_iban", IsUKIBAN)
	validator.RegisterValidation("reference_mask", IsReferenceMask)

	return validator
}
