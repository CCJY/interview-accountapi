package accounts

import (
	"context"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type CreateAccountInterface interface {
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	//
	// When uses this NewCreateAccountRequest, it can be used to RequestInterface that
	// includes WithContext, WithRetry, WhenBeforeDo, Do and WhenAfterDo.
	//
	// To send an HTTP request and return an HTTP response, call Do function.
	NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse]
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	//
	// When uses this CreateAccount, it returns CreateAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This CreateAccount has operation which Client.Do to send an HTTP request and an HTTP response.
	CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	//
	// When uses this CreateAccountWithContext, it returns CreateAccountResponseContext that
	// has http.Response and the ContextData which is user-defined data.
	//
	// This CreateAccountWithContext has operation which Client.Do to
	//  send an HTTP request and an HTTP response.
	// If want to specific context, it can be used.
	CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
}

func (a *AccountClient) NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse] {
	return client.NewRequestContext[types.CreateAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodPost),
			client.WithBody(createAccountRequest),
			client.WithUrl(a.Client.BaseUrl, OperationPathCreateAccount),
		),
	)
}

func (a *AccountClient) CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return a.NewCreateAccountRequest(createAccountRequest).Do()
}

func (a *AccountClient) CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return a.NewCreateAccountRequest(createAccountRequest).
		WithContext(ctx).
		Do()
}
