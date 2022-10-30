package accounts

import (
	"context"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type GetAccountInterface interface {
	// Fetch a single Account resource using the resource ID.
	NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse]
	GetAccount(accountId string) (*types.GetAccountResponseContext, error)
	GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error)
}

func (a *AccountClient) NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse] {
	return client.NewRequestContext[types.GetAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodGet),
			client.WithUrl(a.Client.Config.BaseUrl, OperationPathGetAccount),
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
