package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/ccjy/interview-accountapi/pkg/types"

	account_types "github.com/ccjy/interview-accountapi/pkg/types/account/v1"
)

type HttpClient struct {
	EndpointURL *url.URL
	Client      *http.Client
	Encoding    client.Encoding
}

func New(rawUrl string) (client.ClientInterface, error) {
	c := &http.Client{}

	serverURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	e := &client.JSONEncoding{}

	return &HttpClient{
		EndpointURL: serverURL,
		Client:      c,
		Encoding:    e,
	}, nil
}

// Create a new bank account or register an existing bank account with Form3.
// Since FPS requires accounts to be in the UK, the value of the country attribute must be GB.
func (c *HttpClient) CreateAccountWithResponse(body *account_types.CreateAccountBody) (*account_types.CreateAccountWithResponse, error) {
	reader, err := c.Encoding.Marshal(body)
	if err != nil {
		return nil, err
	}

	queryURL, err := c.EndpointURL.Parse(account_types.OperationPathCreateAccount)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), reader)
	if err != nil {
		return nil, err
	}

	fmt.Println(queryURL.String())
	req.Header.Add("Content-Type", account_types.ContentType)

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var rspData account_types.CreateAccountWithResponse

	err = c.Encoding.UnMarshal(rsp.Body, &rspData)
	if err != nil {
		return nil, err
	}

	rspData.HttpResponse = types.HttpResponse{
		HttpResponse: rsp,
	}

	return &rspData, nil
}

// Fetch a single Account resource using the resource ID.
func (c *HttpClient) GetAccountByIdWithResponse(account_id string) (*account_types.GetAccountByIdWithResponse, error) {
	path := strings.ReplaceAll(account_types.OperationPathGetAccount, "{account_id}", account_id)
	url, err := c.EndpointURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var rspData account_types.GetAccountByIdWithResponse

	err = c.Encoding.UnMarshal(rsp.Body, &rspData)
	if err != nil {
		return nil, err
	}

	rspData.HttpResponse = types.HttpResponse{
		HttpResponse: rsp,
	}

	return &rspData, nil
}

// List accounts with the ability to filter and paginate.
// All accounts that match all filter criteria will be returned (combinations of filters act as AND expressions).
// Multiple values can be set for filters in CSV format, e.g. filter[country]=GB,FR,DE.
func (c *HttpClient) GetAccountAllWithResponse(page *account_types.FilterPage, filters ...*account_types.Filter) (*account_types.GetAccountAllWithResponses, error) {
	url, err := c.EndpointURL.Parse(account_types.OperationPathAllAccount)
	if err != nil {
		return nil, err
	}

	queryValues := url.Query()

	for _, v := range filters {
		queryValues.Add(string(v.Key), string(v.Value))

	}
	if page != nil {
		queryValues.Add("number", fmt.Sprint(page.Number))
		queryValues.Add("size", fmt.Sprint(page.Size))
	}

	url.RawQuery = queryValues.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var rspData account_types.GetAccountAllWithResponses

	err = c.Encoding.UnMarshal(rsp.Body, &rspData)
	if err != nil {
		return nil, err
	}

	rspData.HttpResponse = types.HttpResponse{
		HttpResponse: rsp,
	}

	return &rspData, nil
}

// Delete an Account resource using the resource ID and the current version number.
func (c *HttpClient) DeleteAccountByIdAndVersionWithResponse(account_id string, params *account_types.DeleteAccountByIdAndVersionParams) (*account_types.DeleteAccountByWithResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params must not be nil")
	}
	path := strings.ReplaceAll(account_types.OperationPathDeleteAccount, "{account_id}", account_id)
	url, err := c.EndpointURL.Parse(path)
	if err != nil {
		return nil, err
	}

	queryValues := url.Query()

	queryValues.Add("version", fmt.Sprint(params.Version))

	url.RawQuery = queryValues.Encode()

	req, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var rspData account_types.DeleteAccountByWithResponse

	switch rsp.StatusCode {
	case http.StatusNoContent:
	case http.StatusNotFound:
	case http.StatusConflict:
		break
	case http.StatusBadRequest:
		err = c.Encoding.UnMarshal(rsp.Body, &rspData)
		if err != nil {
			return nil, err
		}
	}

	rspData.HttpResponse = types.HttpResponse{
		HttpResponse: rsp,
	}

	return &rspData, nil
}
