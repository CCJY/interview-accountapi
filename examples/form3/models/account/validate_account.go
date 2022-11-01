package account

import "github.com/go-playground/validator/v10"

func (v *AccountValidator) AccountValidation(sl validator.StructLevel) {
	account := sl.Current().Interface().(AccountData)

	if v.validate.Var(account.Attributes.Iban, "iban") != nil {
		sl.ReportError(account.Attributes.Iban, "iban", "attributes", "iban", "param")
	}

}
