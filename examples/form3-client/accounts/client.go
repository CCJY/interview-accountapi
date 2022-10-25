package accounts

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ccjy/interview-accountapi/examples/form3-client/accounts/types"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type AccountClientInterface interface {
	CreateAccount(account *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error)
	GetAccount(accountId string) (*types.GetAccountResponseContext, error)
	DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error)
}

type AccountClient struct {
	Client *client.Client
}

func New(client *client.Client) AccountClientInterface {
	return &AccountClient{
		Client: client,
	}
}

func (a *AccountClient) CreateAccount(account *types.CreateAccountRequest) (*types.CreateAccountResponseContext, error) {
	return client.NewRequest(
		a.Client,
		&types.CreateAccountRequestContext{
			Method: http.MethodPost,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: OperationPathCreateAccount,
			},
			Body: account,
		},
	).WhenBeforeDo(func(rc *types.CreateAccountRequestContext) error {
		return nil
	}).WhenAfterDo(func(rc *types.CreateAccountResponseContext) error {
		switch rc.StatusCode() {
		case http.StatusCreated:
			fmt.Printf("Created")
		case http.StatusConflict:
		case http.StatusBadRequest:
		default:
			fmt.Printf("unexpected error: %d", rc.StatusCode())
		}
		return nil
	}).Do()
}

func (a *AccountClient) GetAccount(accountId string) (*types.GetAccountResponseContext, error) {
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
	).WhenBeforeDo(func(rc *types.GetAccountRequestContext) error {
		return nil
	}).WhenAfterDo(func(rc *types.GetAccountResponseContext) error {
		switch rc.StatusCode() {
		case http.StatusOK:
			fmt.Printf("get account")
		case http.StatusConflict:
		case http.StatusBadRequest:
		default:
			fmt.Printf("unexpected error: %d", rc.StatusCode())
		}
		return nil
	}).Do()
}

func (a *AccountClient) DeleteAccount(accountId string, version string) (*types.DeleteAccountResponseContext, error) {
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
	).WhenBeforeDo(func(rc *types.DeleteAccountRequestContext) error {
		return nil
	}).WhenAfterDo(func(rc *types.DeleteAccountResponseContext) error {
		switch rc.StatusCode() {
		case http.StatusNoContent:
			fmt.Printf("deleted account")
		case http.StatusConflict:
		case http.StatusBadRequest:
		default:
			fmt.Printf("unexpected error: %d", rc.StatusCode())
		}
		return nil
	}).Do()
}
