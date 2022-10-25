package accounts

import (
	"fmt"
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3-fake-account-api/types"
	"github.com/ccjy/interview-accountapi/examples/form3-fake-account-api/types/account"
	account_types "github.com/ccjy/interview-accountapi/examples/form3-fake-account-api/types/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type CreateAccountRequest types.RequestData[account.AccountData]
type CreateAccountResponse types.ResponseData[account.AccountData]

type ClientDo interface {
	CreateAccount(account *CreateAccountRequest) (*client.ResponseContext[CreateAccountResponse], error)
}

type Client struct {
	Client *client.Client
}

func New(client *client.Client) ClientDo {
	return &Client{
		Client: client,
	}
}

func (a *Client) CreateAccount(account *CreateAccountRequest) (*client.ResponseContext[CreateAccountResponse], error) {
	return client.NewRequest(
		a.Client,
		&client.RequestContext[CreateAccountResponse]{
			Method: http.MethodPost,
			UrlBuilder: &client.Url{
				BaseUrl:       a.Client.Config.BaseUrl,
				OperationPath: account_types.OperationPathCreateAccount,
			},
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: account,
		},
	).WhenBeforeDo(func(rc *client.RequestContext[CreateAccountResponse]) error {
		return nil
	}).WhenAfterDo(func(rc *client.ResponseContext[CreateAccountResponse]) error {
		switch rc.StatusCode() {
		case http.StatusCreated:
			fmt.Printf("Created")
		case http.StatusConflict:
		case http.StatusBadRequest:
		default:
			fmt.Printf("unexpected error")
		}
		return nil
	}).Do()
}
