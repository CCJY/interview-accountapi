package client

import (
	account_types "github.com/ccjy/interview-accountapi/pkg/types/account/v1"
)

type ClientInterface interface {
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	CreateAccountWithResponse(body *account_types.CreateAccountBody) (*account_types.CreateAccountWithResponse, error)

	// Fetch a single Account resource using the resource ID.
	GetAccountByIdWithResponse(account_id string) (*account_types.GetAccountByIdWithResponse, error)

	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	GetAccountAllWithResponse(page *account_types.FilterPage, filter ...*account_types.Filter) (*account_types.GetAccountAllWithResponses, error)

	// Delete an Account resource using the resource ID and the current version number.
	DeleteAccountByIdAndVersionWithResponse(account_id string, params *account_types.DeleteAccountByIdAndVersionParams) (*account_types.DeleteAccountByWithResponse, error)
}
