package form3fakeaccountapi

import (
	"os"

	"github.com/ccjy/interview-accountapi/examples/form3-fake-account-api/accounts"
	"github.com/ccjy/interview-accountapi/pkg/client"
)

type Form3Client struct {
	AccountApi accounts.ClientDo
}

func New() *Form3Client {
	f := &Form3Client{}

	transport := client.NewTransport()

	accountClient := client.NewClient(transport, client.ClientConfig{
		BaseUrl: getHostNmae(),
		Timeout: 3,
	}, nil)

	f.AccountApi = accounts.New(accountClient)
	return f
}

func getHostNmae() string {
	env := os.Getenv("APP-ENV")
	switch env {
	case "docker":
		return "http://accountapi:8080"
	default:
		return "http://127.0.0.1:8080"
	}
}

func (f *Form3Client) GetAccountApi() accounts.ClientDo {
	return f.AccountApi
}
