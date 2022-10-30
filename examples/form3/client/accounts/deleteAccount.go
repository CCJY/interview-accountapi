package accounts

import (
	"context"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type DeleteAccountInterface interface {
	// Delete an Account resource using the resource ID and the current version number.
	//
	// When uses this NewDeleteAccountRequest, it can be used to RequestInterface that
	// includes WithContext, WithRetry, WhenBeforeDo, Do and WhenAfterDo.
	//
	// To send an HTTP request and return an HTTP response, call Do function.
	NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse]
	// Delete an Account resource using the resource ID and the current version number.
	//
	// When uses this DeleteAccount, it returns DeleteAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This DeleteAccount has operation which Client.Do to send an HTTP request and an HTTP response.
	DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error)
	// Delete an Account resource using the resource ID and the current version number.
	//
	// When uses this DeleteAccountWithContext, it returns DeleteAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This DeleteAccountWithContext has operation which Client.Do to
	// send an HTTP request and an HTTP response.
	// If want to specific context, it can be used.
	DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error)
}

func (a *AccountClient) NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse] {
	return client.NewRequestContext[types.DeleteAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodDelete),
			client.WithUrl(a.Client.BaseUrl, OperationPathDeleteAccount),
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
