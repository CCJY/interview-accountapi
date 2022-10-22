package v1

import (
	"github.com/ccjy/interview-accountapi/pkg/types"
)

type RequestData struct {
	Data *AccountData `json:"data,omitempty"`
}

type ResponseData struct {
	types.HttpResponse
	Data         *AccountData `json:"data,omitempty"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

type ResponseDataArray struct {
	types.HttpResponse
	Data         *[]AccountData `json:"data,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
}

type CreateAccountBody RequestData
type CreateAccountWithResponse ResponseData

type GetAccountByIdWithResponse ResponseData

// DeleteAccountByIdAndVersionParams defines parameters for DeleteAccountByIdAndVersion.
type DeleteAccountByIdAndVersionParams struct {
	// Current version number of the Account resource.
	Version int64 `form:"version" json:"version"`
}

type DeleteAccountByWithResponse struct {
	types.HttpResponse
	ErrorMessage string `json:"error_message,omitempty"`
}

// Filters defines model for Filters.
// type Filters struct {
// 	AccountNumber *string `json:"account_number,omitempty"`
// 	BankId        *string `json:"bank_id,omitempty"`
// 	BankIdCode    *string `json:"bank_id_code,omitempty"`
// 	Country       *string `json:"country,omitempty"`
// 	CustomerId    *string `json:"customer_id,omitempty"`
// 	Iban          *string `json:"iban,omitempty"`
// }

// PageFilter defines model for PageFilter.
type FilterPage struct {
	// Page number being requested, defaults to 0.
	Number int

	// Size of the page being requested, defaults to 100.
	Size int
}

// // GetAccountAllParams defines parameters for GetAccountAll.
// type GetAccountAllParams struct {
// 	// Options for filtering the results
// 	Filter *Filters `json:"filter,omitempty"`

// 	// Options for filtering the results
// 	Page *PageFilter `json:"page,omitempty"`
// }

type GetAccountAllWithResponses ResponseDataArray

type FilterKey string

const (
	FilterAccountNumber FilterKey = "account_number"
	FilterBankId        FilterKey = "bank_id"
	FilterBankIdCode    FilterKey = "bank_id_code"
	FilterCountry       FilterKey = "country"
	FilterCustomerId    FilterKey = "customer_id"
	FilterIban          FilterKey = "iban"
)

type Filter struct {
	Key   FilterKey
	Value string
}
