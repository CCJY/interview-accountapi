package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ccjy/interview-accountapi/examples/form3/client/accounts/types"
	"github.com/ccjy/interview-accountapi/examples/form3/models/account"
	"github.com/ccjy/interview-accountapi/pkg/client"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func getHostNmae() string {
	env := os.Getenv("APP-ENV")
	switch env {
	case "docker":
		return "http://host.docker.internal:8080"
	default:
		return "http://127.0.0.1:8080"
	}
}

func getAccountClient() AccountClientInterface {
	transport := client.InitTransport()
	ccjyclient := client.NewClient(
		client.WithTransport(transport),
		client.WithBaseUrl(getHostNmae()),
		client.WithTimeout(3))

	return New(ccjyclient)
}

func getAccountClientWithMockAndTimeout(mockurl string, milliseconds int) AccountClientInterface {
	transport := client.InitTransport()
	ccjyclient := client.NewClient(
		client.WithTransport(transport),
		client.WithBaseUrl(mockurl),
		client.WithTimeout(milliseconds),
	)

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
		Id:             uuid.New().String(),
		OrganisationId: uuid.New().String(),
		Version:        lo.ToPtr(int64(0)),
		Type:           lo.ToPtr("accounts"),
	}
}

func TestAccountClient_CreateAccount_When_CreateAccount_Then_201(t *testing.T) {
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
			Data: account_data,
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

}

func TestAccountClient_CreateAccount_Given_Account_When_ExistsAccount_Then_409(t *testing.T) {
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
			Data: account_data,
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
		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i), http.StatusCreated))
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode(), http.StatusCreated) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.StatusCode(), tt.want)
			}

			tt.name = fmt.Sprintf("expected %d", http.StatusConflict)
			tt.want = http.StatusConflict
			t.Run(tt.name, func(t *testing.T) {
				got1, err := client.CreateAccount(tt.args.account)
				if (err != nil) != tt.wantErr {
					t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got1.StatusCode(), http.StatusConflict) {
					t.Errorf("AccountClient.CreateAccount() = %v, want %v", got1.StatusCode(), tt.want)
				}
			})

		})

	}

}

func TestAccountClient_CreateAccount_When_CreateAccount_Then_201_And_ValidData(t *testing.T) {
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

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("#%d", i)))
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ContextData.Data, tt.want.Data) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.ContextData.Data.Id, tt.want.Data.Id)
			}
		})
	}
}

func TestAccountClient_CreateAccount_Given_InvalidData_When_CreateAccount_Then_400(t *testing.T) {
	type args struct {
		createAccountRequest *types.CreateAccountRequest
	}
	client := getAccountClient()
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "When Data Nil 400",
			args: args{
				createAccountRequest: &types.CreateAccountRequest{
					Data: nil,
				},
			},
			want:    http.StatusBadRequest,
			wantErr: false,
		},
		{
			name: "When Id Emtpy 400",
			args: args{
				createAccountRequest: &types.CreateAccountRequest{
					Data: &account.AccountData{
						Attributes: &account.AccountAttributes{
							BankId:     "400300",
							BankIdCode: "GBDSC",
							Bic:        "NWBKGB22",
							Country:    "GB",
							Name: &[]string{
								"it", "is", "example",
							},
						},
						Id:             "",
						OrganisationId: uuid.New().String(),
						Version:        lo.ToPtr(int64(0)),
						Type:           lo.ToPtr("accounts"),
					},
				},
			},
			want:    http.StatusBadRequest,
			wantErr: false,
		},
		{
			name: "When Invalid Data 400",
			args: args{
				createAccountRequest: &types.CreateAccountRequest{
					Data: &account.AccountData{
						Attributes: &account.AccountAttributes{
							BankId:     "400300",
							BankIdCode: "GBDSC",
							Bic:        "NWBKGB22",
							Country:    "GB",
							// Name: &[]string{
							// 	"it", "is", "example",
							// },
						},
						Id:             uuid.New().String(),
						OrganisationId: uuid.New().String(),
						Version:        lo.ToPtr(int64(0)),
						Type:           lo.ToPtr("accounts"),
					},
				},
			},
			want:    http.StatusBadRequest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateAccount(tt.args.createAccountRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode(), tt.want) {
				t.Errorf("AccountClient.CreateAccount() = %v, want %v", got.StatusCode(), tt.want)
			}
		})
	}
}

func TestAccountClient_GetAccount_When_CreateAccount_Then_200_And_ValidData(t *testing.T) {
	t.Parallel()
	type args struct {
		id string
	}

	type test struct {
		name    string
		args    args
		want    *types.CreateAccountRequest
		wantErr bool
	}

	account_data := DefaultAccountData()
	want := &types.CreateAccountRequest{
		Data: account_data,
	}

	tt := test{
		name: "get account",
		args: args{
			id: want.Data.Id,
		},
		want: want,
	}

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			bytes, _ := json.Marshal(want)
			fmt.Fprintln(w, string(bytes))
		}),
	)

	defer s.Close()

	client := getAccountClientWithMockAndTimeout(s.URL, 300)

	t.Run(tt.name, func(t *testing.T) {
		got1, err := client.GetAccount(tt.args.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got1.ContextData.Data, tt.want.Data) {
			t.Errorf("AccountClient.CreateAccount() = %v, want %v", got1.ContextData.Data, tt.want.Data)
		}
	})
}

func TestAccountClient_DeleteAccount_Given_Account_When_Delete_Then_204(t *testing.T) {
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

	for i := 0; i < 1; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("create #%d", i)))
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
		n := fmt.Sprintf("%s_then_delete", tt.name)
		t.Run(n, func(t *testing.T) {
			got2, err := client.DeleteAccount(tt.args.account.Data.Id, fmt.Sprint(*tt.args.account.Data.Version))
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got2.StatusCode() != http.StatusNoContent {
				t.Errorf("AccountClient.DeleteAccount() = %v, want %v", got2.StatusCode(), http.StatusNoContent)
			}
		})
	}

}

func TestAccountClient_DeleteAccount_When_DeleteAccountNotExists_Then_404(t *testing.T) {
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

	for i := 0; i < 1; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("#%d", i)))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.DeleteAccount(tt.args.account.Data.Id, fmt.Sprint(*tt.args.account.Data.Version))
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode() != http.StatusNotFound {
				t.Errorf("AccountClient.DeleteAccount() = %v, want %v", got.StatusCode(), http.StatusNotFound)
			}
		})
	}

}

func TestAccountClient_CreateAccount_TimeoutDeadline(t *testing.T) {
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

	tests := []test{}

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
			want:    data,
			wantErr: true,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i)))
	}

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(150 * time.Millisecond)
		}),
	)

	defer s.Close()

	client := getAccountClientWithMockAndTimeout(s.URL, 100)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := client.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}

}

func TestAccountClient_CreateAccount_WithContextDeadline(t *testing.T) {
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
			Data: account_data,
		}

		return test{
			name: name,
			args: args{
				data,
			},
			want:    data,
			wantErr: true,
		}
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, createAcountFn(fmt.Sprintf("create account #%d", i)))
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

			defer cancel()

			time.Sleep(150 * time.Millisecond)
			_, err := client.CreateAccountWithContext(ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountClient.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
