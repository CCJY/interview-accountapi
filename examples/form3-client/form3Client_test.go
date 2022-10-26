package form3client

import (
	"reflect"
	"testing"

	"github.com/ccjy/interview-accountapi/examples/form3-client/accounts/types"
	"github.com/ccjy/interview-accountapi/examples/form3-client/models/account"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func DefaultAccountData() *account.AccountData {
	attributes := account.AccountAttributes{
		BankId:     "400300",
		BankIdCode: "GBDSC",
		Bic:        "NWBKGB22",
		Country:    "GB",
		Name: &[]string{
			"it", "is", "example",
		},
	}

	return &account.AccountData{
		Attributes:     &attributes,
		Id:             lo.ToPtr(uuid.New()),
		OrganisationId: lo.ToPtr(uuid.New()),
		Version:        lo.ToPtr(int64(0)),
		Type:           lo.ToPtr("accounts"),
	}
}

func TestForm3Client_CreateAccount(t *testing.T) {
	type args struct {
		accountData *types.CreateAccountRequest
	}
	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	createAcountFn := func(name string) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: data,
		}
	}
	tests := []test{
		// TODO: Add test cases.
		createAcountFn("Create Account"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New()
			got, err := a.GetAccountApi().CreateAccount(tt.args.accountData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.want.Data, got.ContextData.Data) {
				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// assert.Equal(t, tt.args.accountData.Data.Id, got.ContextData.Data.Id)
			// assert.Equal(t, tt.args.accountData.Data.OrganisationId, got.ContextData.Data.OrganisationId)
			// assert.Equal(t, tt.args.accountData.Data.Version, got.ContextData.Data.Version)
			// assert.Equal(t, tt.args.accountData.Data.Type, got.ContextData.Data.Type)
			// assert.Equal(t, tt.args.accountData.Data.Attributes.BankId, got.ContextData.Data.Attributes.BankId)
			// assert.Equal(t, tt.args.accountData.Data.Attributes.BankIdCode, got.ContextData.Data.Attributes.BankIdCode)
			// assert.Equal(t, tt.args.accountData.Data.Attributes.Bic, got.ContextData.Data.Attributes.Bic)
			// assert.Equal(t, tt.args.accountData.Data.Attributes.Country, got.ContextData.Data.Attributes.Country)
			// assert.Equal(t, tt.args.accountData.Data.Attributes.Name, got.ContextData.Data.Attributes.Name)
		})
	}

}

// func TestForm3Client_GetAccountApi(t *testing.T) {
// 	type fields struct {
// 		AccountApi accounts.AccountClientInterface
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   accounts.AccountClientInterface
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f := &Form3Client{
// 				AccountApi: tt.fields.AccountApi,
// 			}
// 			if got := f.GetAccountApi(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Form3Client.GetAccountApi() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func DefaultRequestData() *account_types.CreateAccountBody {
// 	return &account_types.CreateAccountBody{
// 		Data: DefaultAccountData(),
// 	}
// }

// func AssertCreateAccountWithResponse(t *testing.T, statusCode int, request *account_types.CreateAccountBody, resp *account_types.CreateAccountWithResponse) {
// 	switch statusCode {
// 	case http.StatusCreated:
// 		assert.Equal(t, request.Data.Id, resp.Data.Id)
// 		assert.Equal(t, request.Data.OrganisationId, resp.Data.OrganisationId)
// 		assert.Equal(t, request.Data.Version, resp.Data.Version)
// 		assert.Equal(t, request.Data.Type, resp.Data.Type)
// 		assert.Equal(t, request.Data.Attributes.BankId, resp.Data.Attributes.BankId)
// 		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.Data.Attributes.BankIdCode)
// 		assert.Equal(t, request.Data.Attributes.Bic, resp.Data.Attributes.Bic)
// 		assert.Equal(t, request.Data.Attributes.Country, resp.Data.Attributes.Country)
// 		assert.Equal(t, request.Data.Attributes.Name, resp.Data.Attributes.Name)

// 	case http.StatusConflict:
// 	case http.StatusBadRequest:
// 		t.Errorf("CreateAccountWithResponse() status = %d, err = %s", statusCode, resp.ErrorMessage)
// 	}
// }

// func AssertGetAccountByIdWithResponse(t *testing.T, statusCode int, request *account_types.CreateAccountBody, resp *account_types.GetAccountByIdWithResponse) {
// 	switch statusCode {
// 	case http.StatusOK:
// 		assert.Equal(t, request.Data.Id, resp.Data.Id)
// 		assert.Equal(t, request.Data.OrganisationId, resp.Data.OrganisationId)
// 		assert.Equal(t, request.Data.Version, resp.Data.Version)
// 		assert.Equal(t, request.Data.Type, resp.Data.Type)
// 		assert.Equal(t, request.Data.Attributes.BankId, resp.Data.Attributes.BankId)
// 		assert.Equal(t, request.Data.Attributes.BankIdCode, resp.Data.Attributes.BankIdCode)
// 		assert.Equal(t, request.Data.Attributes.Bic, resp.Data.Attributes.Bic)
// 		assert.Equal(t, request.Data.Attributes.Country, resp.Data.Attributes.Country)
// 		assert.Equal(t, request.Data.Attributes.Name, resp.Data.Attributes.Name)

// 	case http.StatusNotFound:
// 		t.Errorf("GetAccountByIdWithResponse() status = %d, err = %s", statusCode, resp.ErrorMessage)
// 	}
// }

// func AssertDeleteAccountByIdAndVersionWithResponse(t *testing.T, statusCode int, request *account_types.CreateAccountBody, resp *account_types.DeleteAccountByWithResponse) {
// 	switch statusCode {
// 	case http.StatusNoContent:
// 	case http.StatusNotFound:
// 	case http.StatusConflict:
// 		return
// 	case http.StatusBadRequest:
// 		t.Errorf("DeleteAccountByIdAndVersionWithResponse() status = %d, err = %s", statusCode, resp.ErrorMessage)
// 	}
// }

// func TestPost(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	rsp, err := client.NewRequestContext(&client.RequestContext[account_types.CreateAccountWithResponse]{

// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathCreateAccount,
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 		Method: http.MethodPost,
// 		Body:   reqData,
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertCreateAccountWithResponse(t, rsp.StatusCode(), reqData, &rsp.Data)
// }

// func TestPostWithoutEncoding(t *testing.T) {

// 	client.Init(client.WithNewTransport(10))
// 	reqData := DefaultRequestData()

// 	rsp, err := client.NewRequestContext(&client.RequestContext[account_types.CreateAccountWithResponse]{
// 		Method:        http.MethodPost,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathCreateAccount,
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 		Body: reqData,
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertCreateAccountWithResponse(t, rsp.StatusCode(), reqData, &rsp.Data)
// }

// func TestPostWithoutInit(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	rsp, err := client.NewRequestContext(&client.RequestContext[account_types.CreateAccountWithResponse]{
// 		Method:        http.MethodPost,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathCreateAccount,
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 		Body: reqData,
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertCreateAccountWithResponse(t, rsp.StatusCode(), reqData, &rsp.Data)
// }

// func TestRequestContext_DeadlineExceeded(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

// 	defer cancel()

// 	requestContext := client.NewRequestContext(&client.RequestContext[account_types.CreateAccountWithResponse]{
// 		Method:        http.MethodPost,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathCreateAccount,
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 		Body:    reqData,
// 		Context: ctx,
// 	})

// 	time.Sleep(2 * time.Second)

// 	_, err := requestContext.Do()

// 	assert.Error(t, err)

// }

// func TestGetAccountByIdWithResponse_NoExists404(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	version := fmt.Sprint(reqData.Data.Version)

// 	rsp, err := client.NewRequestContext(&client.RequestContext[account_types.GetAccountByIdWithResponse]{
// 		Method:        http.MethodGet,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathGetAccount,
// 		PathParams: map[string]string{
// 			"account_id": reqData.Data.Id.String(),
// 		},
// 		QueryParams: url.Values{
// 			"version": []string{version},
// 		},
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, http.StatusNotFound, rsp.StatusCode())

// }

// func TestGetAccountByIdWithResponse_200(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	rspCreate, err := client.NewRequestContext(&client.RequestContext[account_types.CreateAccountWithResponse]{
// 		Method:        http.MethodPost,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathCreateAccount,
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 		Body: reqData,
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertCreateAccountWithResponse(t, rspCreate.StatusCode(), reqData, &rspCreate.Data)

// 	version := fmt.Sprint(reqData.Data.Version)

// 	rspGet, err := client.NewRequestContext(&client.RequestContext[account_types.GetAccountByIdWithResponse]{
// 		Method:        http.MethodGet,
// 		BaseUrl:       baseUrl,
// 		OperationPath: account_types.OperationPathGetAccount,
// 		PathParams: map[string]string{
// 			"account_id": reqData.Data.Id.String(),
// 		},
// 		QueryParams: url.Values{
// 			"version": []string{version},
// 		},
// 		Header: map[string][]string{
// 			"Content-Type": {"application/json"},
// 		},
// 	}).Do()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertGetAccountByIdWithResponse(t, rspGet.StatusCode(), reqData, &rspGet.Data)
// }

// func TestNewRequestContext_CreateAccountWithResponse_201(t *testing.T) {
// 	reqData := DefaultRequestData()

// 	requestCreateContext := client.NewRequestContext(
// 		&client.RequestContext[account_types.CreateAccountWithResponse]{
// 			Method:        http.MethodPost,
// 			BaseUrl:       baseUrl,
// 			OperationPath: account_types.OperationPathCreateAccount,
// 			Header: map[string][]string{
// 				"Content-Type": {"application/json"},
// 			},
// 			Body: reqData,
// 		})

// 	rspCreate, err := requestCreateContext.Do()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	AssertCreateAccountWithResponse(t, rspCreate.StatusCode(), reqData, &rspCreate.Data)
// }

// func TestNewRequestContext_CreateAccountWithResponse_Multipe3(t *testing.T) {
// 	reqData1 := DefaultRequestData()
// 	reqData2 := DefaultRequestData()
// 	reqData3 := DefaultRequestData()

// 	var waitGroup sync.WaitGroup
// 	c := make(chan error)
// 	waitGroup.Add(3)

// 	rspDatas := struct {
// 		rspCreate1 *client.ResponseContext[account_types.CreateAccountWithResponse]
// 		rspCreate2 *client.ResponseContext[account_types.CreateAccountWithResponse]
// 		rspCreate3 *client.ResponseContext[account_types.CreateAccountWithResponse]
// 	}{}

// 	go func() {
// 		fmt.Printf("1 Waiting...")
// 		waitGroup.Wait()
// 		fmt.Printf("Closing...")
// 		close(c)
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp1, err := client.NewRequestContext(
// 			&client.RequestContext[account_types.CreateAccountWithResponse]{
// 				Method:        http.MethodPost,
// 				BaseUrl:       baseUrl,
// 				OperationPath: account_types.OperationPathCreateAccount,
// 				Header: map[string][]string{
// 					"Content-Type": {"application/json"},
// 				},
// 				Body: reqData1,
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate1 = rsp1
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp2, err := client.NewRequestContext(
// 			&client.RequestContext[account_types.CreateAccountWithResponse]{
// 				Method:        http.MethodPost,
// 				BaseUrl:       baseUrl,
// 				OperationPath: account_types.OperationPathCreateAccount,
// 				Header: map[string][]string{
// 					"Content-Type": {"application/json"},
// 				},
// 				Body: reqData2,
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate2 = rsp2
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp3, err := client.NewRequestContext(
// 			&client.RequestContext[account_types.CreateAccountWithResponse]{
// 				Method:        http.MethodPost,
// 				BaseUrl:       baseUrl,
// 				OperationPath: account_types.OperationPathCreateAccount,
// 				Body:          reqData3,
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate3 = rsp3
// 	}()

// 	fmt.Printf("2 Waiting...")
// 	for err := range c {
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	fmt.Printf("3 Waiting...")

// 	AssertCreateAccountWithResponse(t, rspDatas.rspCreate1.StatusCode(), reqData1, &rspDatas.rspCreate1.Data)
// 	AssertCreateAccountWithResponse(t, rspDatas.rspCreate2.StatusCode(), reqData2, &rspDatas.rspCreate2.Data)
// 	AssertCreateAccountWithResponse(t, rspDatas.rspCreate3.StatusCode(), reqData3, &rspDatas.rspCreate3.Data)
// }

// // https://httpbin.org/get
// type TestModel struct {
// 	Url    string `json:"url"`
// 	Origin string `json:"origin"`
// }

// func TestNewRequestContext_DifferentUrl(t *testing.T) {
// 	reqData1 := DefaultRequestData()
// 	reqData2 := DefaultRequestData()

// 	var waitGroup sync.WaitGroup
// 	c := make(chan error)
// 	waitGroup.Add(3)

// 	rspDatas := struct {
// 		rspCreate1 *client.ResponseContext[account_types.CreateAccountWithResponse]
// 		rspCreate2 *client.ResponseContext[account_types.CreateAccountWithResponse]
// 		rspCreate3 *client.ResponseContext[TestModel]
// 	}{}

// 	go func() {
// 		fmt.Printf("1 Waiting...")
// 		waitGroup.Wait()
// 		fmt.Printf("Closing...")
// 		close(c)
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp1, err := client.NewRequestContext(
// 			&client.RequestContext[account_types.CreateAccountWithResponse]{
// 				Method:        http.MethodPost,
// 				BaseUrl:       baseUrl,
// 				OperationPath: account_types.OperationPathCreateAccount,
// 				Header: map[string][]string{
// 					"Content-Type": {"application/json"},
// 				},
// 				Body: reqData1,
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate1 = rsp1
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp2, err := client.NewRequestContext(
// 			&client.RequestContext[account_types.CreateAccountWithResponse]{
// 				Method:        http.MethodPost,
// 				BaseUrl:       baseUrl,
// 				OperationPath: account_types.OperationPathCreateAccount,
// 				Header: map[string][]string{
// 					"Content-Type": {"application/json"},
// 				},
// 				Body: reqData2,
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate2 = rsp2
// 	}()

// 	go func() {
// 		defer waitGroup.Done()
// 		rsp3, err := client.NewRequestContext(
// 			&client.RequestContext[TestModel]{
// 				Method:        http.MethodGet,
// 				BaseUrl:       "https://httpbin.org",
// 				OperationPath: "get",
// 			}).Do()
// 		if err != nil {
// 			c <- err
// 		}
// 		rspDatas.rspCreate3 = rsp3
// 	}()

// 	fmt.Printf("2 Waiting...")
// 	for err := range c {
// 		if err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	fmt.Printf("3 Waiting...")

// 	AssertCreateAccountWithResponse(t, rspDatas.rspCreate1.StatusCode(), reqData1, &rspDatas.rspCreate1.Data)
// 	AssertCreateAccountWithResponse(t, rspDatas.rspCreate2.StatusCode(), reqData2, &rspDatas.rspCreate2.Data)

// 	assert.NotNil(t, rspDatas.rspCreate3.Data.Url)
// }
