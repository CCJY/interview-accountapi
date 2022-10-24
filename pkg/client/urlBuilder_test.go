package client

import (
	"net/url"
	"reflect"
	"testing"
)

func TestUrl_Build(t *testing.T) {
	type fields struct {
		BaseUrl       string
		OperationPath string
		QueryParams   url.Values
		PathParams    map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				BaseUrl:       "http://127.0.0.1:8080",
				OperationPath: "/v1/organisation/accounts",
				PathParams: map[string]string{
					"account_id": "wefoiaejf",
				},
				QueryParams: url.Values{
					"filter[account_id]": []string{"account_id"},
				},
			},
			want: "http://127.0.0.1:8080/v1/organisation/accounts?filter[account_id]=account_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Url{
				BaseUrl:       tt.fields.BaseUrl,
				OperationPath: tt.fields.OperationPath,
				QueryParams:   tt.fields.QueryParams,
				PathParams:    tt.fields.PathParams,
			}
			got, err := u.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Url.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			queryUnescape, err := url.QueryUnescape(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Url.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(queryUnescape, tt.want) {
				t.Errorf("Url.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
