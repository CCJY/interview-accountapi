package accounts

import (
	"context"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type GetAccountInterface interface {
	// Fetch a single Account resource using the resource ID.
	//
	// When uses this NewGetAccountRequest, it can be used to RequestInterface that
	// includes WithContext, WithRetry, WhenBeforeDo, Do and WhenAfterDo.
	//
	// To send an HTTP request and return an HTTP response, call Do function.
	NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse]
	// Fetch a single Account resource using the resource ID.
	//
	// When uses this GetAccount, it returns GetAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This GetAccount has operation which Client.Do to
	// send an HTTP request and an HTTP response.
	GetAccount(accountId string) (*types.GetAccountResponseContext, error)
	// Fetch a single Account resource using the resource ID.
	//
	// When uses this GetAccountWithContext, it returns GetAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This GetAccountWithContext has operation which Client.Do to
	// send an HTTP request and an HTTP response.
	// If want to specific context, it can be used.
	GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error)
}

func (a *AccountClient) NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse] {
	return client.NewRequestContext[types.GetAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodGet),
			client.WithUrl(a.Client.BaseUrl, OperationPathGetAccount),
			client.WithPathParams(client.WithPathParam("account_id", accountId)),
		),
	)
}

func (a *AccountClient) GetAccount(accountId string) (*types.GetAccountResponseContext, error) {
	return a.NewGetAccountRequest(accountId).Do()
}

func (a *AccountClient) GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error) {
	return a.NewGetAccountRequest(accountId).
		WithContext(ctx).
		Do()
}
