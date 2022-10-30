package accounts

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type GetAllAccountInterface interface {
	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	//
	// When uses this NewGetAllAccountRequest, it can be used to RequestInterface that
	// includes WithContext, WithRetry, WhenBeforeDo, Do and WhenAfterDo.
	//
	// To send an HTTP request and return an HTTP response, call Do function.
	NewGetAllAccountRequest(opts ...types.GetAllAccountOpt) client.RequestInterface[types.GetAllAccountResponse]
	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	//
	// When uses this GetAllAccount, it returns GetAllAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This GetAllAccount has operation which Client.Do to
	// send an HTTP request and an HTTP response.
	GetAllAccount(opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error)
	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	//
	// When uses this GetAllAccountWithContext, it returns GetAllAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This GetAllAccountWithContext has operation which Client.Do to
	// send an HTTP request and an HTTP response.
	// If want to specific context, it can be used.
	GetAllAccountWithContext(ctx context.Context, opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error)
}

func (a *AccountClient) NewGetAllAccountRequest(opts ...types.GetAllAccountOpt) client.RequestInterface[types.GetAllAccountResponse] {
	queryValues := url.Values{}

	for _, opt := range opts {
		opt(&queryValues)
	}

	return client.NewRequestContext[types.GetAllAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodGet),
			client.WithUrl(a.Client.BaseUrl, OperationPathAllAccount),
			client.WithQueryValues(&queryValues),
		),
	)
}

func (a *AccountClient) GetAllAccount(opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error) {
	return a.NewGetAllAccountRequest(opts...).Do()
}
func (a *AccountClient) GetAllAccountWithContext(ctx context.Context, opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error) {
	return a.NewGetAllAccountRequest(opts...).
		WithContext(ctx).
		Do()
}
