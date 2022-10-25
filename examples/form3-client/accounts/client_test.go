package accounts

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/ccjy/interview-accountapi/examples/form3-client/accounts/types"
	"github.com/ccjy/interview-accountapi/examples/form3-client/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func getHostNmae() string {
	env := os.Getenv("APP-ENV")
	switch env {
	case "docker":
		return "http://accountapi:8080"
	default:
		return "http://127.0.0.1:8080"
	}
}

func getAccountClient() AccountClientInterface {
	transport := client.NewTransport()
	ccjyclient := client.NewClient(transport, client.ClientConfig{
		BaseUrl: getHostNmae(),
		Timeout: 3,
	}, nil)

	return New(ccjyclient)
}

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

func TestAccountClient_CreateAccount_StatusCode(t *testing.T) {
	t.Parallel()
	type args struct {
		account *types.CreateAccountRequest
	}

	type test struct {
		name    string
		args    args
		want    int
		wantErr bool
	}

	client := getAccountClient()

	tests := []test{}

	createAcountFn := func(name string, statusCode int) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: *account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: statusCode,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("#%d", i), http.StatusCreated))
	}

	for _, tt := range tests {
		tt.name = fmt.Sprintf("%s, status: %d", tt.name, http.StatusCreated)
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode(), http.StatusCreated) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want)
			}
		})

	}

	for _, tt := range tests {
		tt.name = fmt.Sprintf("%s, status: %d", tt.name, http.StatusConflict)
		tt.want = http.StatusConflict
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode(), http.StatusConflict) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want)
			}
		})

	}
}

func TestAccountClient_CreateAccount_ValidResponseData(t *testing.T) {
	t.Parallel()
	type args struct {
		account *types.CreateAccountRequest
	}

	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	client := getAccountClient()

	tests := []test{}

	createAcountFn := func(name string) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: *account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: data,
		}
	}

	for i := 0; i < 1; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("#%d", i)))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}
}

func TestAccountClient_GetAccount_ValidResponseData(t *testing.T) {
	t.Parallel()
	type args struct {
		account *types.CreateAccountRequest
	}

	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	client := getAccountClient()

	tests := []test{}

	createAcountFn := func(name string) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: *account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: data,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i)))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}

	for _, tt := range tests {
		n := fmt.Sprintf("get account against %s", tt.name)
		t.Run(n, func(t *testing.T) {
			got, err := client.GetAccount(tt.args.account.Data.Id.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}
}

func TestAccountClient_DeleteAccount_204(t *testing.T) {
	t.Parallel()
	type args struct {
		account *types.CreateAccountRequest
	}

	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	client := getAccountClient()

	tests := []test{}

	createAcountFn := func(name string) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: *account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: data,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("#%d", i)))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}

	for _, tt := range tests {
		n := fmt.Sprintf("get account against create account %s", tt.name)
		t.Run(n, func(t *testing.T) {
			got, err := client.GetAccount(tt.args.account.Data.Id.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.GetAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}

	for _, tt := range tests {
		n := fmt.Sprintf("delete account against create account %s", tt.name)
		t.Run(n, func(t *testing.T) {
			got, err := client.DeleteAccount(tt.args.account.Data.Id.String(), fmt.Sprint(*tt.args.account.Data.Version))
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode() != http.StatusNoContent {
				t.Errorf("AccountClient.DeleteAccount() = %v, want %v", got.StatusCode(), http.StatusNoContent)
			}
		})
	}

}

func TestAccountClient_GetAccount_Timeout(t *testing.T) {
	t.Parallel()
	type args struct {
		account *types.CreateAccountRequest
	}

	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	client := getAccountClient()

	tests := []test{}

	createAcountFn := func(name string) test {
		account_data := DefaultAccountData()
		data := &types.CreateAccountRequest{
			Data: *account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want: data,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i)))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}

	for _, tt := range tests {
		n := fmt.Sprintf("get account against : %s", tt.name)
		t.Run(n, func(t *testing.T) {
			got, err := client.GetAccount(tt.args.account.Data.Id.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
			}
		})
	}
}

// func TestAccountClient_GetAccount_WithContext(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		account *types.CreateAccountRequest
// 	}

// 	type test struct {
// 		name    string
// 		args    args
// 		want    *types.CreateAccountRequest
// 		wantErr bool
// 	}

// 	client := getAccountClient()

// 	tests := []test{}

// 	createAcountFn := func(name string) test {
// 		account_data := DefaultAccountData()
// 		data := &types.CreateAccountRequest{
// 			Data: *account_data,
// 		}

// 		return test{
// 			name: name,
// 			args: args{
// 				data,
// 			},
// 			want: data,
// 		}
// 	}

// 	for i := 0; i < 10; i++ {
// 		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i)))
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := client.CreateAccount(tt.args.account)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
// 				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
// 			}
// 		})
// 	}

// 	for _, tt := range tests {
// 		n := fmt.Sprintf("get account against : %s", tt.name)
// 		t.Run(n, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

// 			defer cancel()

// 			got, err := client.GetAccount(tt.args.account.Data.Id.String())
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
// 				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data, tt.want.Data)
// 			}
// 		})
// 	}
// }
