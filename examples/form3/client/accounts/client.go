package accounts

import (
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type AccountClientInterface interface {
	CreateAccountInterface
	DeleteAccountInterface
	GetAccountInterface
	GetAllAccountInterface
	HealthInterface
}

type AccountClient struct {
	Client *client.Client
}

func New(client *client.Client) AccountClientInterface {
	return &AccountClient{
		Client: client,
	}
}
