package accounts

import (
	"context"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type DeleteAccountInterface interface {
	// Delete an Account resource using the resource ID and the current version number.
	NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse]
	DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error)
	DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error)
}

func (a *AccountClient) NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse] {
	return client.NewRequestContext[types.DeleteAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodDelete),
			client.WithUrl(a.Client.Config.BaseUrl, OperationPathDeleteAccount),
			client.WithPathParams(client.WithPathParam("account_id", accountId)),
			client.WithQueryParams(client.WithQueryParam("version", version)),
		),
	)
}

func (a *AccountClient) DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error) {
	return a.NewDeleteAccountRequest(accountId, version).Do()
}
func (a *AccountClient) DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error) {
	return a.NewDeleteAccountRequest(accountId, version).
		WithContext(ctx).
		Do()
}
