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
	NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse]
	CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
}

func (a *AccountClient) NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse] {
	return client.NewRequestContext[types.CreateAccountResponse](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodPost),
			client.WithBody(createAccountRequest),
			client.WithUrl(a.Client.Config.BaseUrl, OperationPathCreateAccount),
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
