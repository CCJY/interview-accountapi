package types

import (
	"fmt"
	"net/url"

	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	models "github.com/ccjy/interview-accountapi/examples/form3/models/account"
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
	Key   FilterKey `json:"key,omitempty"`
	Value string    `json:"value,omitempty"`
}

type GetAllAccountParams struct {
	Page    *FilterPage `json:"page,omitempty"`
	Filters []Filter    `json:"filters,omitempty"`
}

type GetAllAccountOpt func(*url.Values)

func WithPage(number int, size int) GetAllAccountOpt {

	return func(urlValues *url.Values) {
		urlValues.Set("page[number]", fmt.Sprint(number))
		urlValues.Set("page[size]", fmt.Sprint(size))
	}
}

func WithFilter(key string, value string) GetAllAccountOpt {
	return func(urlValues *url.Values) {
		filter_key := fmt.Sprintf("filter[%s]", key)
		urlValues.Add(filter_key, fmt.Sprint(value))
	}
}

type GetAllAccountResponse = commons.ResponseDataArray[models.AccountData]
type GetAllAccountResponseContext = client.ResponseContext[GetAllAccountResponse]
