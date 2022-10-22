package http

import (
	"fmt"
	"net/http"
	"testing"

	account_types "github.com/ccjy/interview-accountapi/pkg/types/account"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var (
	serverUrl = "http://127.0.0.1:8080/v1/"
)

func DefaultAccountData() *account_types.AccountData {
	attributes := account_types.AccountAttributes{
		BankId:     "400300",
		BankIdCode: "GBDSC",
		Bic:        "NWBKGB22",
		Country:    "GB",
		Name: &[]string{
			"it", "is", "example",
		},
	}

	return &account_types.AccountData{
		Attributes:     &attributes,
		Id:             lo.ToPtr(uuid.New()),
		OrganisationId: lo.ToPtr(uuid.New()),
		Version:        lo.ToPtr(int64(0)),
		Type:           lo.ToPtr("accounts"),
	}
}

func DefaultRequestData() *account_types.CreateAccountBody {
	return &account_types.CreateAccountBody{
		Data: DefaultAccountData(),
	}
}

func AssertCreateAccountWithResponse(t *testing.T, request *account_types.CreateAccountBody, resp *account_types.CreateAccountWithResponse) {
	switch resp.StatusCode() {
	case http.StatusCreated:
		assert.Equal(t, request.Data.Id, resp.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, resp.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, resp.Data.Version)
		assert.Equal(t, request.Data.Type, resp.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, resp.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, resp.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, resp.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, resp.Data.Attributes.Name)

	case http.StatusBadRequest:
		t.Errorf("CreateAccountWithResponse() status = %d, err = %s", resp.StatusCode(), resp.ErrorMessage)
	}
}

func AssertGetAccountByIdWithResponse(t *testing.T, request *account_types.CreateAccountBody, resp *account_types.GetAccountByIdWithResponse) {

	switch resp.StatusCode() {
	case http.StatusOK:
		assert.Equal(t, request.Data.Id, resp.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, resp.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, resp.Data.Version)
		assert.Equal(t, request.Data.Type, resp.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, resp.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, resp.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, resp.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, resp.Data.Attributes.Name)

	case http.StatusNotFound:
		t.Errorf("GetAccountByIdWithResponse() status = %d, err = %s", resp.StatusCode(), resp.ErrorMessage)
	}
}

func AssertDeleteAccountByIdAndVersionWithResponse(t *testing.T, request *account_types.CreateAccountBody, resp *account_types.DeleteAccountByWithResponse) {
	switch resp.StatusCode() {
	case http.StatusNoContent:
	case http.StatusNotFound:
	case http.StatusConflict:
		return
	case http.StatusBadRequest:
		t.Errorf("DeleteAccountByIdAndVersionWithResponse() status = %d, err = %s", resp.StatusCode(), resp.ErrorMessage)
	}
}

func TestHttpClient_CreateAccountWithResponse(t *testing.T) {

	type args struct {
		body *account_types.CreateAccountBody
	}
	tests := []struct {
		name    string
		args    args
		want    *account_types.CreateAccountWithResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				body: DefaultRequestData(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(serverUrl)
			if err != nil {
				t.Errorf("HttpClient.CreateAccountWithResponse() error = %v", err)
				return
			}

			if got, err := c.CreateAccountWithResponse(tt.args.body); err != nil {
				t.Errorf("HttpClient.CreateAccountWithResponse() error = %v", err)
			} else {
				AssertCreateAccountWithResponse(t, tt.args.body, got)
			}
		})
	}
}

func TestCreateAndGetAndDeleteAccount(t *testing.T) {
	requestBody := DefaultRequestData()

	client, err := New(serverUrl)
	if err != nil {
		t.Errorf("HttpClient.CreateAccountWithResponse() error = %v", err)
		return
	}

	if got, err := client.CreateAccountWithResponse(requestBody); err != nil {
		t.Errorf("%v", err)
	} else {
		AssertCreateAccountWithResponse(t, requestBody, got)
	}

	if got, err := client.GetAccountByIdWithResponse(requestBody.Data.Id.String()); err != nil {
		t.Errorf("%v", err)
	} else {
		AssertGetAccountByIdWithResponse(t, requestBody, got)
	}

	params := account_types.DeleteAccountByIdAndVersionParams{Version: *requestBody.Data.Version}
	if got, err := client.DeleteAccountByIdAndVersionWithResponse(requestBody.Data.Id.String(), &params); err != nil {
		fmt.Println(err)
	} else {
		AssertDeleteAccountByIdAndVersionWithResponse(t, requestBody, got)
	}

}

func TestCreate10AndGetAllAndDeleteAllAccount(t *testing.T) {
	client, err := New(serverUrl)
	if err != nil {
		t.Errorf("HttpClient.CreateAccountWithResponse() error = %v", err)
		return
	}
	for i := 0; i < 10; i++ {
		request := DefaultRequestData()

		if got, err := client.CreateAccountWithResponse(request); err != nil {
			t.Errorf("%v", err)
		} else {
			AssertCreateAccountWithResponse(t, request, got)
		}
	}

	accounts := make(map[string]int64)

	page := account_types.FilterPage{
		Number: 1,
		Size:   10,
	}

	accountFilter := &account_types.Filter{
		Key:   account_types.FilterAccountNumber,
		Value: "",
	}

	if got, err := client.GetAccountAllWithResponse(&page, accountFilter); err != nil {
		t.Errorf("%v", err)
	} else {
		assert.Equal(t, http.StatusOK, got.StatusCode(), "Should be: %v, got: %v", http.StatusOK, got.StatusCode())
		if got.Data != nil {
			for _, data := range *got.Data {
				accounts[data.Id.String()] = *data.Version
			}
		}
	}

	for key, data := range accounts {
		params := account_types.DeleteAccountByIdAndVersionParams{Version: data}
		if got, err := client.DeleteAccountByIdAndVersionWithResponse(key, &params); err != nil {
			t.Errorf("%v", err)
		} else {
			assert.Equal(t, http.StatusNoContent, got.StatusCode(), "Should be: %v, got: %v", http.StatusNoContent, got.StatusCode())
		}
	}

	if got, err := client.GetAccountAllWithResponse(&account_types.FilterPage{}); err != nil {
		t.Errorf("%v", err)
	} else {
		assert.Equal(t, http.StatusOK, got.StatusCode(), "Should be: %v, got: %v", http.StatusOK, got.StatusCode())
		assert.Nil(t, got.Data)
	}
}
