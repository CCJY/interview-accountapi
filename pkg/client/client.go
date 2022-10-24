package client

import (
	types "github.com/ccjy/interview-accountapi/pkg/client/types/account"
)

type ClientInterface interface {
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	CreateAccountWithResponse(body *types.CreateAccountBody) (*types.CreateAccountWithResponse, error)

	// Fetch a single Account resource using the resource ID.
	GetAccountByIdWithResponse(account_id string) (*types.GetAccountByIdWithResponse, error)

	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	GetAccountAllWithResponse(page *types.FilterPage, filter ...*types.Filter) (*types.GetAccountAllWithResponses, error)

	// Delete an Account resource using the resource ID and the current version number.
	DeleteAccountByIdAndVersionWithResponse(account_id string, params *types.DeleteAccountByIdAndVersionParams) (*types.DeleteAccountByWithResponse, error)
}
