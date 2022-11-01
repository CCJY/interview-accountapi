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

	_ "net/http/pprof"
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

func getTransport() *client.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 30
	t.MaxConnsPerHost = 30
	t.MaxIdleConnsPerHost = 30
	t.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).DialContext

	transport := client.InitTransport(
		client.WithNewTransport(t),
	)

	return transport
}

func main() {
	t := getTransport()
	serverUrl := "localhost:9090"
	client := client.NewClient(
		client.WithTransport(t),
		client.WithBaseUrl("http://"+serverUrl),
		client.WithTimeout(100),
	)

	handlers := wrapperStruct{client: client}
	http.HandleFunc("/req", handlers.requestHandler)
	http.HandleFunc("/wait", handlers.waitHandler)
	http.ListenAndServe(serverUrl, nil)
}

type wrapperStruct struct {
	client *client.Client
}

func (ws wrapperStruct) requestHandler(w http.ResponseWriter, r *http.Request) {
	reqData := &types.CreateAccountRequest{
		Data: getDefaultAccountData(),
	}

	got, err := client.NewRequestContext[types.CreateAccountResponse](
		ws.client,
		client.NewRequestContextModel(
			client.WithHttpMethod(http.MethodPost),
			client.WithBody(reqData),
			client.WithUrl(ws.client.BaseUrl, "/wait"),
		),
	).WithRetry(
		client.WithRetryPolicyExpoFullyBackOff(100, 300, 3),
	).Do()

	if err != nil {
		fmt.Printf("%v", err)
		fmt.Fprintln(w, err)
		return
	}

	dataBytes, err := json.Marshal(got.ContextData)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Fprintln(w, dataBytes)
		return
	}

	fmt.Fprintln(w, string(dataBytes))
}

func (ws wrapperStruct) waitHandler(w http.ResponseWriter, r *http.Request) {
	reqData := new(types.CreateAccountRequest)
	err := json.NewDecoder(r.Body).Decode(reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}

	dataBytes, err := json.Marshal(reqData)
	if err != nil {
		w.Write(dataBytes)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(dataBytes))
}

func HealthCheck(accountClient accounts.AccountClientInterface) {
	got, err := accountClient.HealthCheck()

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("status: %s\n", got.ContextData.Status)
}
