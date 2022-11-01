package account

import (
	"encoding/json"
	"fmt"
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

func IsReferenceMask(fl validator.FieldLevel) bool {
	regex := "^[#$?\\\\]{1,35}?$"
	ReferenceMaskRegex := regexp.MustCompile(regex)

	return ReferenceMaskRegex.MatchString(fl.Field().String())
}

type AccountValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewAccountValidator() *AccountValidator {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	v := &AccountValidator{
		validate: validator.New(),
		trans:    trans,
	}
	v.validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(uuid.UUID); ok {
			return valuer.String()
		}
		return nil
	}, uuid.UUID{})

	// CustomerRegister(validator, trans, "swift_code", "", IsSWiftCode)
	CustomerRegister(v.validate, trans, "bank_account_number", "{0}{1}", IsBankAccountNumber)
	CustomerRegister(v.validate, trans, "iban", "{0} {1}", IsIBAN)
	// CustomerRegister(validator, trans, "reference_mask", "", IsReferenceMask)

	v.validate.RegisterStructValidation(v.AccountValidation, &AccountData{})
	return v
}

func CustomerRegister(v *validator.Validate, trans ut.Translator, tag string, text string, validFn func(f1 validator.FieldLevel) bool) {
	v.RegisterValidation(tag, validFn)
	v.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			return fe.(error).Error()
		}
		return t
	})

}
func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
