package accounts

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type AccountClientInterface interface {
	NewCreateAccountRequest(createAccountRequest *types.CreateAccountRequest) client.RequestInterface[types.CreateAccountResponse]
	NewGetAccountRequest(accountId string) client.RequestInterface[types.GetAccountResponse]
	NewGetAllAccountRequest() client.RequestInterface[types.GetAllAccountResponse]
	NewDeleteAccountRequest(accountId string, version string) client.RequestInterface[types.DeleteAccountResponse]

	CreateAccount(createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	GetAccount(accountId string) (*types.GetAccountResponseContext, error)
	GetAllAccount() (*types.GetAllAccountResponseContext, error)
	DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error)

	CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error)
	GetAllAccountWithContext(ctx context.Context) (*types.GetAllAccountResponseContext, error)
	DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error)
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

func (a *AccountClient) NewGetAllAccountRequest() client.RequestInterface[types.GetAllAccountResponse] {
	return client.NewRequest(
		a.Client,
		&types.GetAllAccountRequestContext{
			Method: http.MethodGet,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathAllAccount,
			},
		},
	)
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

func (a *AccountClient) GetAccount(accountId string) (*types.GetAccountResponseContext, error) {
	return a.NewGetAccountRequest(accountId).Do()
}

func (a *AccountClient) GetAllAccount() (*types.GetAllAccountResponseContext, error) {
	return a.NewGetAllAccountRequest().Do()
}

func (a *AccountClient) DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error) {
	return a.NewDeleteAccountRequest(accountId, version).Do()
}

func (a *AccountClient) CreateAccountWithContext(ctx context.Context, createAccountRequest *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return a.NewCreateAccountRequest(createAccountRequest).
		WithContext(ctx).
		Do()
}

func (a *AccountClient) GetAllAccountWithContext(ctx context.Context) (*types.GetAllAccountResponseContext, error) {
	return a.NewGetAllAccountRequest().
		WithContext(ctx).
		Do()
}

func (a *AccountClient) GetAccountWithContext(ctx context.Context, accountId string) (*types.GetAccountResponseContext, error) {
	return a.NewGetAccountRequest(accountId).
		WithContext(ctx).
		Do()
}

func (a *AccountClient) DeleteAccountWithContext(ctx context.Context, accountId string, version string) (*types.DeleteAccountResponseContext, error) {
	return a.NewDeleteAccountRequest(accountId, version).
		WithContext(ctx).
		Do()
}
