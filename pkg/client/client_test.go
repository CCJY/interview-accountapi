package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

var (
	serverUrl = "http://127.0.0.1:8080/v1/"
)

func TestCreateAndGetAndDeleteAccount(t *testing.T) {
	attributes := AccountAttributes{
		BankId:     "400300",
		BankIdCode: "GBDSC",
		Bic:        "NWBKGB22",
		Country:    "GB",
		Name: &[]string{
			"it", "is", "example",
		},
	}

	request := RequestData{
		Data: &AccountData{
			Attributes:     &attributes,
			Id:             lo.ToPtr(uuid.New()),
			OrganisationId: lo.ToPtr(uuid.New()),
			Version:        lo.ToPtr(int64(0)),
			Type:           lo.ToPtr("accounts"),
		},
	}

	client := New(serverUrl)

	if got, err := client.CreateAccountWithResponse(serverUrl, request); err != nil {
		t.Errorf("%v", err)
	} else {
		assert.Equal(t, http.StatusCreated, got.StatusCode(), "Should be: %v, got: %v", http.StatusCreated, got.StatusCode())
		assert.Equal(t, request.Data.Id, got.JSON201.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, got.JSON201.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, got.JSON201.Data.Version)
		assert.Equal(t, request.Data.Type, got.JSON201.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, got.JSON201.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, got.JSON201.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, got.JSON201.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, got.JSON201.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, got.JSON201.Data.Attributes.Name)
	}

	if got, err := client.GetAccountByIdWithResponse(serverUrl, request.Data.Id.String()); err != nil {
		t.Errorf("%v", err)
	} else {
		assert.Equal(t, http.StatusOK, got.StatusCode(), "Should be: %v, got: %v", http.StatusOK, got.StatusCode())
		assert.Equal(t, request.Data.Id, got.JSON200.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, got.JSON200.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, got.JSON200.Data.Version)
		assert.Equal(t, request.Data.Type, got.JSON200.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, got.JSON200.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, got.JSON200.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, got.JSON200.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, got.JSON200.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, got.JSON200.Data.Attributes.Name)
	}

	params := DeleteAccountByIdAndVersionParams{Version: *request.Data.Version}
	got, err := client.DeleteAccountByIdAndVersionWithResponse(serverUrl, request.Data.Id.String(), &params)

	if err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusNoContent, got.StatusCode())
	assert.Equal(t, []byte{}, got.HTTPResponse.Body())
}
func TestCreate10AndGetAllAndDeleteAllAccount(t *testing.T) {
	client := New(serverUrl)
	CreateAccountsTest(t, client, 10)

	GetAllAndDeleteAllAccountTest(t, client)
}

func CreateAccountsTest(t *testing.T, client *Client, loopLength int) {
	for i := 0; i < loopLength; i++ {
		attributes := AccountAttributes{
			BankId:     "400300",
			BankIdCode: "GBDSC",
			Bic:        "NWBKGB22",
			Country:    "GB",
			Name: &[]string{
				"it", "is", "example",
			},
		}

		request := RequestData{
			Data: &AccountData{
				Attributes:     &attributes,
				Id:             lo.ToPtr(uuid.New()),
				OrganisationId: lo.ToPtr(uuid.New()),
				Version:        lo.ToPtr(int64(0)),
				Type:           lo.ToPtr("accounts"),
			},
		}

		if got, err := client.CreateAccountWithResponse(serverUrl, request); err != nil {
			t.Errorf("%v", err)
		} else {
			assert.Equal(t, http.StatusCreated, got.StatusCode(), "Should be: %v, got: %v", http.StatusCreated, got.StatusCode())
			assert.Equal(t, request.Data.Id, got.JSON201.Data.Id)
			assert.Equal(t, request.Data.OrganisationId, got.JSON201.Data.OrganisationId)
			assert.Equal(t, request.Data.Version, got.JSON201.Data.Version)
			assert.Equal(t, request.Data.Type, got.JSON201.Data.Type)
			assert.Equal(t, request.Data.Attributes.BankId, got.JSON201.Data.Attributes.BankId)
			assert.Equal(t, request.Data.Attributes.BankIdCode, got.JSON201.Data.Attributes.BankIdCode)
			assert.Equal(t, request.Data.Attributes.Bic, got.JSON201.Data.Attributes.Bic)
			assert.Equal(t, request.Data.Attributes.Country, got.JSON201.Data.Attributes.Country)
			assert.Equal(t, request.Data.Attributes.Name, got.JSON201.Data.Attributes.Name)
		}
	}
}

func GetAllAndDeleteAllAccountTest(t *testing.T, client *Client) {

	params := GetAccountAllParams{}
	accounts := make(map[string]int64)

	if got, err := client.GetAccountAllWithResponse(serverUrl, &params); err != nil {
		t.Errorf("%v", err)
	} else {
		assert.Equal(t, http.StatusOK, got.StatusCode(), "Should be: %v, got: %v", http.StatusOK, got.StatusCode())
		if got.JSON200.Data != nil {
			for _, data := range *got.JSON200.Data {
				accounts[data.Id.String()] = *data.Version
			}
		}
	}

	for key, data := range accounts {
		params := DeleteAccountByIdAndVersionParams{Version: data}
		if got, err := client.DeleteAccountByIdAndVersionWithResponse(serverUrl, key, &params); err != nil {
			t.Errorf("%v", err)
		} else {
			assert.Equal(t, http.StatusNoContent, got.StatusCode(), "Should be: %v, got: %v", http.StatusNoContent, got.StatusCode())
		}
	}

}

// func GetAllAccountWithFilterCountryTest(t *testing.T, client *Client, expectedLength int, expectedCountry string) {
// 	for i := 0; i < expectedLength; i++ {
// 		attributes := AccountAttributes{
// 			BankId:     "400300",
// 			BankIdCode: "GBDSC",
// 			Bic:        "NWBKGB22",
// 			Country:    expectedCountry,
// 			Name: &[]string{
// 				"it", "is", "example",
// 			},
// 		}

// 		request := RequestData{
// 			Data: &AccountData{
// 				Attributes:     &attributes,
// 				Id:             lo.ToPtr(uuid.New()),
// 				OrganisationId: lo.ToPtr(uuid.New()),
// 				Version:        lo.ToPtr(int64(0)),
// 				Type:           lo.ToPtr("accounts"),
// 			},
// 		}

// 		if got, err := client.CreateAccountWithResponse(serverUrl, request); err != nil {
// 			t.Errorf("%v", err)
// 		} else {
// 			assert.Equal(t, http.StatusCreated, got.StatusCode(), "Should be: %v, got: %v", http.StatusCreated, got.StatusCode())
// 			assert.Equal(t, request.Data.Id, got.JSON201.Data.Id)
// 			assert.Equal(t, request.Data.OrganisationId, got.JSON201.Data.OrganisationId)
// 			assert.Equal(t, request.Data.Version, got.JSON201.Data.Version)
// 			assert.Equal(t, request.Data.Type, got.JSON201.Data.Type)
// 			assert.Equal(t, request.Data.Attributes.BankId, got.JSON201.Data.Attributes.BankId)
// 			assert.Equal(t, request.Data.Attributes.BankIdCode, got.JSON201.Data.Attributes.BankIdCode)
// 			assert.Equal(t, request.Data.Attributes.Bic, got.JSON201.Data.Attributes.Bic)
// 			assert.Equal(t, request.Data.Attributes.Country, got.JSON201.Data.Attributes.Country)
// 			assert.Equal(t, request.Data.Attributes.Name, got.JSON201.Data.Attributes.Name)
// 		}
// 	}

// 	params := GetAccountAllParams{Filter: &Filters{
// 		Country: lo.ToPtr(expectedCountry),
// 	}, Page: &PageFilter{
// 		Number: lo.ToPtr("1"),
// 		Size:   lo.ToPtr("10"),
// 	}}
// 	if got, err := client.GetAccountAllWithResponse(serverUrl, &params); err != nil {
// 		t.Errorf("%v", err)
// 	} else {
// 		assert.Equal(t, http.StatusOK, got.StatusCode(), "Should be: %v, got: %v", http.StatusOK, got.StatusCode())
// 		dataLength := len(*got.JSON200.Data)
// 		assert.Equal(t, expectedLength, dataLength)
// 		for _, data := range *got.JSON200.Data {
// 			assert.Equal(t, expectedCountry, data.Attributes.Country)
// 		}
// 	}
// }
