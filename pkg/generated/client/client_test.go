package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var (
	serverUrl = "http://127.0.0.1:8080/v1/"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func DefaultAccountData() *AccountData {
	attributes := AccountAttributes{
		BankId:     "400300",
		BankIdCode: "GBDSC",
		Bic:        "NWBKGB22",
		Country:    "GB",
		Name: &[]string{
			"it", "is", "example",
		},
	}

	return &AccountData{
		Attributes:     &attributes,
		Id:             lo.ToPtr(uuid.New()),
		OrganisationId: lo.ToPtr(uuid.New()),
		Version:        lo.ToPtr(int64(0)),
		Type:           lo.ToPtr("accounts"),
	}
}

func DefaultRequestData() *RequestData {

	return &RequestData{
		Data: DefaultAccountData(),
	}
}

func AssertCreateAccountWithResponse(t *testing.T, request *RequestData, resp *CreateAccountResponse) {

	switch resp.StatusCode() {
	case http.StatusCreated:
		assert.Equal(t, request.Data.Id, resp.JSON201.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, resp.JSON201.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, resp.JSON201.Data.Version)
		assert.Equal(t, request.Data.Type, resp.JSON201.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, resp.JSON201.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.JSON201.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, resp.JSON201.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, resp.JSON201.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, resp.JSON201.Data.Attributes.Name)
		assert.NotNil(t, resp.JSON201.Data.CreatedOn)
		assert.NotNil(t, resp.JSON201.Data.ModifiedOn)
		assert.NotNil(t, resp.JSON201.Links)
		assert.NotNil(t, resp.JSON201.Links.Self)
	case http.StatusBadRequest:
		assert.NotNil(t, resp.JSON400.ErrorMessage)
	}
}

func AssertGetAccountByIdWithResponse(t *testing.T, request *RequestData, resp *GetAccountByIdResponse) {

	switch resp.StatusCode() {
	case http.StatusOK:
		assert.Equal(t, request.Data.Id, resp.JSON200.Data.Id)
		assert.Equal(t, request.Data.OrganisationId, resp.JSON200.Data.OrganisationId)
		assert.Equal(t, request.Data.Version, resp.JSON200.Data.Version)
		assert.Equal(t, request.Data.Type, resp.JSON200.Data.Type)
		assert.Equal(t, request.Data.Attributes.BankId, resp.JSON200.Data.Attributes.BankId)
		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.JSON200.Data.Attributes.BankIdCode)
		assert.Equal(t, request.Data.Attributes.Bic, resp.JSON200.Data.Attributes.Bic)
		assert.Equal(t, request.Data.Attributes.Country, resp.JSON200.Data.Attributes.Country)
		assert.Equal(t, request.Data.Attributes.Name, resp.JSON200.Data.Attributes.Name)
		assert.NotNil(t, resp.JSON200.Data.CreatedOn)
		assert.NotNil(t, resp.JSON200.Data.ModifiedOn)
		assert.NotNil(t, resp.JSON200.Links)
		assert.NotNil(t, resp.JSON200.Links.Self)
	case http.StatusNotFound:
		assert.NotNil(t, resp.JSON404.ErrorMessage)
	}
}

func AssertDeleteAccountByIdAndVersionWithResponse(t *testing.T, request *RequestData, resp *DeleteAccountByIdAndVersionResponse) {
	switch resp.StatusCode() {
	case http.StatusNoContent:
	case http.StatusNotFound:
	case http.StatusConflict:
		assert.Equal(t, []byte{}, resp.HTTPResponse.Body())
	case http.StatusBadRequest:
		assert.NotNil(t, resp.JSON400.ErrorMessage)
	}
}

func TestCreateAndGetAndDeleteAccount(t *testing.T) {
	request := DefaultRequestData()

	client := New(serverUrl)

	if got, err := client.CreateAccountWithResponse(serverUrl, *request); err != nil {
		t.Errorf("%v", err)
	} else {
		AssertCreateAccountWithResponse(t, request, got)
	}

	if got, err := client.GetAccountByIdWithResponse(serverUrl, request.Data.Id.String()); err != nil {
		t.Errorf("%v", err)
	} else {
		AssertGetAccountByIdWithResponse(t, request, got)
	}

	params := DeleteAccountByIdAndVersionParams{Version: *request.Data.Version}
	if got, err := client.DeleteAccountByIdAndVersionWithResponse(serverUrl, request.Data.Id.String(), &params); err != nil {
		fmt.Println(err)
	} else {
		AssertDeleteAccountByIdAndVersionWithResponse(t, request, got)
	}

}
func TestCreate10AndGetAllAndDeleteAllAccount(t *testing.T) {
	client := New(serverUrl)
	for i := 0; i < 10; i++ {
		request := DefaultRequestData()

		if got, err := client.CreateAccountWithResponse(serverUrl, *request); err != nil {
			t.Errorf("%v", err)
		} else {
			AssertCreateAccountWithResponse(t, request, got)
		}
	}

	params := GetAccountAllParams{}

	if got, err := client.GetAccountAllWithResponse(serverUrl, &params); err != nil {
		t.Errorf("%v", err)
	} else {
		if got.JSON200.Data != nil {
			for _, data := range *got.JSON200.Data {
				params := DeleteAccountByIdAndVersionParams{Version: *data.Version}
				if got, err := client.DeleteAccountByIdAndVersionWithResponse(serverUrl, data.Id.String(), &params); err != nil {
					t.Errorf("%v", err)
				} else {
					assert.Equal(t, http.StatusNoContent, got.StatusCode(), "Should be: %v, got: %v", http.StatusNoContent, got.StatusCode())
				}
			}
		}
	}

}

func TestValidateParameterRequired(t *testing.T) {
	client := New(serverUrl)
	if got, err := client.GetAccountByIdWithResponse(serverUrl, ""); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(got)
	}

}

func TestValidateDataId(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Id = nil
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Id = lo.ToPtr(uuid.New())
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// organization id required and must be uuid
func TestValidateDataOrganisationId(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.OrganisationId = nil
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.OrganisationId = lo.ToPtr(uuid.New())
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// bank_id requied and max length is 11
func TestValidateDataAttributesBankId(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.BankId = "400300400300400300400300"
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.BankId = "400300"
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// bank_id_code must be GBDSC
func TestValidateDataAttributesBankIdCode(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.BankIdCode = "EOIJFEFE"
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.BankIdCode = "GBDSC"
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesBIC(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Bic = "EOIJFEFEA"
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Bic = "NWBKGB22"
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Bic = "EOIJFEFEA12"
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesCountry(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Country = "GBE"
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Country = "GB"
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesAcceptanceQualifier(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				v := AccountAttributesAcceptanceQualifier(AfterNextWorkingDay)
				accountData.Attributes.AcceptanceQualifier = &v
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesAccountClassification(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				a := AccountAttributesAccountClassification(Business)
				accountData.Attributes.AccountClassification = &a
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesAccountNumber(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AccountNumber = lo.ToPtr("EOIJFEEFE")
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AccountNumber = lo.ToPtr("08464524")
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// alternative names [3][140]
func TestValidateDataAttributesAlternativeNames(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeNames = &[]string{"a", "b", "c", "d"}
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeNames = &[]string{"add", "bdd", "1aaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaef"}
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeNames = &[]string{"add", "bdd", "aa"}
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// currency ISO 4217 must be GBP
func TestValidateDataAttributesCurrency(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				baseCurrency := AccountAttributesBaseCurrency(GBP)
				accountData.Attributes.BaseCurrency = &baseCurrency
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// customer_id is free format max length 36
func TestValidateDataAttributesCustomerId(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				//length 37
				accountData.Attributes.CustomerId = lo.ToPtr("aefaefaefaefaefaefaefaefaefaefaefaeaa")
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				//length 36
				accountData.Attributes.CustomerId = lo.ToPtr("efaefaefaefaefaefaefaefaefaefaefaeB")
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// iban ISO 13616
func TestValidateDataAttributesIban(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Iban = lo.ToPtr("aefaefaefaefaefaefaefaefaefaefaefaeaa")
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Iban = lo.ToPtr("GB11NWBK40030041426811")
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesName(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Name = &[]string{"a", "b", "c", "d", "e"}
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Name = &[]string{"add", "bdd", "1aaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaefaeffaefaefaef"}
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Name = &[]string{"add", "bdd", "aa"}
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.Name = &[]string{"add", "bdd", "aa", "a"}
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesNameMatchingStatus(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				s := AccountAttributesNameMatchingStatus(NotSupported)
				accountData.Attributes.NameMatchingStatus = &s
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesReferenceMask(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.ReferenceMask = lo.ToPtr("########################################################################################################")
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.ReferenceMask = lo.ToPtr("############")
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.ReferenceMask = lo.ToPtr("####$$$\\$######\\#?###\\#\\#\\#######")
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.ReferenceMask = lo.ToPtr("##-#$$$\\$######\\#?###\\#\\#\\#######")
				return accountData
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesSecondaryIdentification(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.SecondaryIdentification = lo.ToPtr(String(141))
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.SecondaryIdentification = lo.ToPtr(String(1))
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesUserDefinedData(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.UserDefinedData = &[]UserDefinedData{}

				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()

				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDataAttributesValidationType(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesValidationType(Card)
				accountData.Attributes.ValidationType = &e

				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()

				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Status of the account. pending and confirmed are set by Form3, closed can be set manually
func TestValidateDataAttributesStatus(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatus(AccountAttributesStatusClosed)
				accountData.Attributes.Status = &e

				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "Error, confirmed is set by Form3",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatus(AccountAttributesStatusConfirmed)
				accountData.Attributes.Status = &e

				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error, confirmed is set by Form3",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatus(AccountAttributesStatusPending)
				accountData.Attributes.Status = &e
				return accountData
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Additional account status information, required when updating status to closed, can't be used otherwise. The FPS code with which inbound payments to the account will be rejected depends on the value of this field.
func TestValidateDataAttributesStatusReason(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error #0",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatusClosed
				accountData.Attributes.Status = &e
				reason := AccountAttributesStatusReasonClosed
				accountData.Attributes.StatusReason = &reason
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "Error, to use it status should be closed #1",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatusConfirmed
				accountData.Attributes.Status = &e
				reason := AccountAttributesStatusReasonClosed
				accountData.Attributes.StatusReason = &reason
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error, to use it status should be closed #2",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				e := AccountAttributesStatusPending
				accountData.Attributes.Status = &e
				reason := AccountAttributesStatusReasonClosed
				accountData.Attributes.StatusReason = &reason
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error, to use it, status should be close and status required #3",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				reason := AccountAttributesStatusReasonClosed
				accountData.Attributes.StatusReason = &reason
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "No Error #4",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				Status := AccountAttributesStatusClosed
				accountData.Attributes.Status = &Status
				return accountData
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// AlternativeBankAccountNames [3][140]
func TestValidateDataAttributesAlternativeBankAccountNames(t *testing.T) {
	tests := []struct {
		name    string
		data    *AccountData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No Error #0",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeBankAccountNames = &[]string{"first", "second", "third"}
				return accountData
			}(),
			wantErr: false,
		},
		{
			name: "Error, #1",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeBankAccountNames = &[]string{"first", "second", "third", "fourth"}
				return accountData
			}(),
			wantErr: true,
		},
		{
			name: "Error, #2",
			data: func() *AccountData {
				accountData := DefaultAccountData()
				accountData.Attributes.AlternativeBankAccountNames = &[]string{"first", String(141)}
				return accountData
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := NewValidator()

			err := validate.Struct(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// func TestReferenceMark(t *testing.T) {
// 	// validate := NewValidator()
// 	// s := "arandomsensitive information: 1234567890 this is not senstive: 1234567890000000"
// 	// re := regexp.MustCompile(`\b(\d{4})\d{6}\b`)
// 	// s = re.ReplaceAllString(s, "$1******$2")
// 	// fmt.Println(s)

// }

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

// https://programming.guide/go/format-parse-string-time-date-example.html
func TestDateTime(t *testing.T) {
	strDate := "2022-10-22T00:33:57.876Z"

	if date, err := time.Parse(time.RFC3339Nano, strDate); err != nil {
		t.Error(err)
	} else {
		fmt.Println(date.String())
	}

}
