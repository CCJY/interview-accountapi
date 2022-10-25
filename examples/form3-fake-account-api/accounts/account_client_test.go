package accounts

// func TestClient_CreateAccount(t *testing.T) {
// 	type fields struct {
// 		Client *client.Client
// 	}
// 	type args struct {
// 		account *account_types.AccountData
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *client.ResponseContext[CreateAccountResponse]
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &Client{
// 				Client: tt.fields.Client,
// 			}
// 			got, err := a.CreateAccount(tt.args.account)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.CreateAccount() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
