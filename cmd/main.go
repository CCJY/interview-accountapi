package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts"
	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/examples/form3/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func getDefaultAccountData() *account.AccountData {
	attributes := account.AccountAttributes{
		BankId:     "400300",
		BankIdCode: "GBDSC",
		Bic:        "NWBKGB22",
		Country:    "GB",
		Name: &[]string{
			"it", "is", "example",
		},
	}

	accountData := &account.AccountData{
		Attributes:     &attributes,
		Id:             uuid.New().String(),
		OrganisationId: uuid.New().String(),
		Version:        lo.ToPtr(int64(0)),
		Type:           lo.ToPtr("accounts"),
	}

	return accountData
}

func main() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 30
	t.MaxConnsPerHost = 30
	t.MaxIdleConnsPerHost = 30
	t.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).DialContext

	transport := client.NewTransport(
		client.WithNewTransport(t),
	)

	client := client.NewClient(transport, client.ClientConfig{
		BaseUrl: "http://127.0.0.1:8080",
	}, nil)

	accountClient := accounts.New(client)

	got, err := accountClient.CreateAccount(&types.CreateAccountRequest{
		Data: getDefaultAccountData(),
	})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	dataBytes, err := json.Marshal(got.ContextData)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("ID: %s\n", got.ContextData.Data.Id)
	fmt.Println("JSON:" + string(dataBytes))
}
