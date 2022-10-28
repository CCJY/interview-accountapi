package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRetry_TRetryRequest(t *testing.T) {
	type fields struct {
		RetryInterval time.Duration
		RetryMax      int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
		want1  bool
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				RetryInterval: 1000,
				RetryMax:      3,
			},
			want:  true,
			want1: true,
		},
		// {
		// 	fields: fields{
		// 		RetryInterval: 200,
		// 		RetryMax:      3,
		// 	},
		// 	want:  false,
		// 	want1: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Retry{
				RetryInterval: tt.fields.RetryInterval,
				RetryMax:      tt.fields.RetryMax,
			}
			got, got1 := r.TRetryRequest()
			if got != tt.want {
				t.Errorf("Retry.RetryRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Retry.RetryRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type TestData struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}
type TestServerResponse struct {
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage *string     `json:"error_message,omitempty"`
}

func TestRetry_RetryRequest(t *testing.T) {
	type fields struct {
		RetryInterval time.Duration
		RetryMax      int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		argsFn func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args   *args
		want   *TestServerResponse
		want1  int //statusc ode
	}{
		// TODO: Add test cases.
		{
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "",
				data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			want: &TestServerResponse{
				Data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
				ErrorMessage: nil,
			},
			want1: 200,
		},

		{
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			want:  nil,
			want1: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retried := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				retried += 1
				var result *TestServerResponse
				if retried < tt.fields.RetryMax {
					w.WriteHeader(500)
					result = tt.want
					dataBytes, _ := json.Marshal(result)

					fmt.Fprintln(w, string(dataBytes))

					return
				}
				w.WriteHeader(tt.want1)
				result = tt.want
				dataBytes, _ := json.Marshal(result)

				fmt.Fprintln(w, string(dataBytes))
			}))

			defer s.Close()

			c := &http.Client{}
			r := &Retry{
				RetryInterval: tt.fields.RetryInterval,
				RetryMax:      tt.fields.RetryMax,
			}
			reader, _ := json.Marshal(tt.args.data)
			request := tt.argsFn(tt.args.method, s.URL, bytes.NewReader(reader))
			got := r.RetryRequest(c, request)

			var responseData TestServerResponse
			bodybytes, err := io.ReadAll(got.Response.Body)
			if err != nil {
				t.Errorf("%v", err)
			}

			defer got.Response.Body.Close()

			err = json.Unmarshal(bodybytes, &responseData)
			if err != nil {
				t.Errorf("%v", err)
			}

			if got.Response.StatusCode != tt.want1 {
				t.Errorf("Retry.RetryRequest() = %v, want1 %v", got.Response.StatusCode, tt.want1)
			}

			var wantData TestServerResponse
			wantBytes, _ := json.Marshal(tt.want)
			err = json.Unmarshal(wantBytes, &wantData)
			if err != nil {
				t.Errorf("%v", err)
			}

			if !reflect.DeepEqual(responseData, wantData) {
				responseDataBytes, _ := json.Marshal(responseData)
				wantBytes, _ := json.Marshal(tt.want)
				t.Errorf("Retry.RetryRequest() = %v, want %v", string(responseDataBytes), string(wantBytes))
			}

		})
	}
}
