package accounts

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type AccountClientInterface interface {
	// Create a new bank account or register an existing bank account with Form3.
	// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
	NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse]
	CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)

	// Delete an Account resource using the resource ID and the current version number.
	NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse]
	DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error)
	DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error)

	// Fetch a single Account resource using the resource ID.
	NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse]
	GetAccount(accountId string) (*types.GetAccountResponseContext, error)
	GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error)

	// List accounts with the ability to filter and paginate.
	// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
	// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
	NewGetAllAccountRequest(opts ...types.GetAllAccountOpt) client.RequestInterface[types.GetAllAccountResponse]
	GetAllAccount(opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error)
	GetAllAccountWithContext(ctx context.Context, opts ...types.GetAllAccountOpt) (*types.GetAllAccountResponseContext, error)
}

type AccountClient struct {
	Client *client.Client
}

func New(client *client.Client) AccountClientInterface {
	return &AccountClient{
		Client: client,
	}
}

func (a *AccountClient) NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse] {
	return client.NewRequest(
		a.Client,
		&types.CreateAccountRequestContext{
			Method: http.MethodPost,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathCreateAccount,
			},
			Body: createAccountRequest,
		},
	)
}

func (a *AccountClient) CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return a.NewCreateAccountRequest(createAccountRequest).
		WhenBeforeDo(func(rc *types.CreateAccountRequestContext) error {
			return nil
		}).WhenAfterDo(func(rc *types.CreateAccountResponseContext) error {
		// switch rc.StatusCode() {
		// case http.StatusCreated:
		// 	fmt.Printf("Created")
		// case http.StatusConflict:
		// case http.StatusBadRequest:
		// default:
		// 	fmt.Printf("unexpected error: %d", rc.StatusCode())
		// }
		return nil
	}).Do()
}

func (a *AccountClient) CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return a.NewCreateAccountRequest(createAccountRequest).
		WithContext(ctx).
		Do()
}

func (a *AccountClient) NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse] {
	return client.NewRequest(
		a.Client,
		&types.DeleteAccountRequestContext{
			Method: http.MethodDelete,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathDeleteAccount,
				PathParams: map[string]string{
					"account_id": accountId,
				},
				QueryParams: url.Values{
					"version": []string{version},
				},
			},
		},
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

func (a *AccountClient) NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse] {
	return client.NewRequest(
		a.Client,
		&types.GetAccountRequestContext{
			Method: http.MethodGet,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathGetAccount,
				PathParams: map[string]string{
					"account_id": accountId,
				},
			},
		},
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

func (a *AccountClient) NewGetAllAccountRequest(opts ...types.GetAllAccountOpt) client.RequestInterface[types.GetAllAccountResponse] {
	queryValues := url.Values{}

	for _, opt := range opts {
		opt(&queryValues)
	}

	return client.NewRequest(
		a.Client,
		&types.GetAllAccountRequestContext{
			Method: http.MethodGet,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathAllAccount,
				QueryParams:   queryValues,
			},
		},
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
