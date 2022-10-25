package types

import (
	"github.com/ccjy/interview-accountapi/examples/form3-client/commons"
	models "github.com/ccjy/interview-accountapi/examples/form3-client/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

// PageFilter defines model for PageFilter.
type FilterPage struct {
	// Page number being requested, defaults to 0.
	Number int

	// Size of the page being requested, defaults to 100.
	Size int
}

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

type GetAccountResponse = commons.ResponseData[models.AccountData]

type GetAccountRequestContext = client.RequestContext[GetAccountResponse]
type GetAccountResponseContext = client.ResponseContext[GetAccountResponse]
