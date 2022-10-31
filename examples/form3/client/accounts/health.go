package accounts

import (
	"net/http"

	"github.com/ccjy/interview-accountapi/examples/form3/commons"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type HealthInterface interface {
	HealthCheck() (*client.ResponseContext[commons.Health], error)
}

func (a *AccountClient) HealthCheck() (*client.ResponseContext[commons.Health], error) {
	return client.NewRequestContext[commons.Health](
		a.Client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodGet),
			client.WithUrl(a.Client.BaseUrl, OperationPathHeatlhCheckAccountAPI),
		),
	).Do()
}
