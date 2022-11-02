package account

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func (v *AccountValidator) AccountValidation(sl validator.StructLevel) {
	account := sl.Current().Interface().(AccountData)

	if v.validate.Var(account.Attributes.Country, "required") != nil {
		sl.ReportError(account.Attributes.Country, "country", "attributes", "required", "")
	}

	if v.validate.Var(account.Attributes.Country, "iso3166_1_alpha2") != nil {
		sl.ReportError(account.Attributes.Country, "country", "attributes", "iso3166_1_alpha2", "")
	}

	validations := MapAttributeValid[account.Attributes.Country]
	if validations == nil {
		return
	}

	for field, validation := range validations {
		for _, tag := range validation.subtags {
			if v.validate.Var(account.Attributes.Iban, tag) != nil {
				sl.ReportError(account.Attributes.Iban, "iban", "attributes", tag, "param")
			}
		}

		tag := fmt.Sprintf("%s_%s", field, account.Attributes.Country)

		switch field {
		case "iban":
			if validation.Regex != nil {
				if v.validate.Var(account.Attributes.Iban, tag) != nil {
					sl.ReportError(account.Attributes.Iban, field, "attributes", tag, validation.Error)
				}
			}
		}

	}

}
