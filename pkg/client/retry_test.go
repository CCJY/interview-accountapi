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
		RetryInterval int
		RetryMax      int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		argsFn          func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args            *args
		serverTimeoutMs int
		clientTimeoutMs int
		want            *TestServerResponse
		wantStatusCode  int //statusc ode
		wantError       bool
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
			wantStatusCode: 200,
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
			want:           nil,
			wantStatusCode: 200,
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
			want:            nil,
			wantStatusCode:  200,
			serverTimeoutMs: 500,
			clientTimeoutMs: 200,
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
			want:            nil,
			wantStatusCode:  200,
			serverTimeoutMs: 500,
			clientTimeoutMs: 200,
			wantError:       true,
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
			want:            nil,
			wantStatusCode:  200,
			serverTimeoutMs: 500,
			clientTimeoutMs: 200,
			wantError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retried := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				retried += 1
				if 0 < tt.serverTimeoutMs && retried < tt.fields.RetryMax || tt.wantError {
					time.Sleep(time.Duration(tt.serverTimeoutMs) * time.Millisecond)
				}
				var result *TestServerResponse
				if retried < tt.fields.RetryMax {
					w.WriteHeader(500)
					result = tt.want
					dataBytes, _ := json.Marshal(result)

					fmt.Fprintln(w, string(dataBytes))

					return
				}
				w.WriteHeader(tt.wantStatusCode)
				result = tt.want
				dataBytes, _ := json.Marshal(result)

				fmt.Fprintln(w, string(dataBytes))
			}))

			defer server.Close()

			c := &http.Client{
				Timeout: time.Duration(tt.clientTimeoutMs) * time.Millisecond,
			}
			r := &Retry{
				RetryInterval: tt.fields.RetryInterval,
				RetryMax:      tt.fields.RetryMax,
			}
			dataBytes, _ := json.Marshal(tt.args.data)
			request := tt.argsFn(tt.args.method, server.URL, bytes.NewReader(dataBytes))
			got, err := r.Do(c, request, dataBytes)

			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}
			if tt.wantError {
				return
			}

			var responseData TestServerResponse
			bodybytes, err := io.ReadAll(got.Body)
			if err != nil {
				t.Errorf("%v", err)
			}

			defer got.Body.Close()

			err = json.Unmarshal(bodybytes, &responseData)
			if err != nil {
				t.Errorf("%v", err)
			}

			if got.StatusCode != tt.wantStatusCode {
				t.Errorf("Retry.RetryRequest() = %v, want1 %v", got.StatusCode, tt.wantStatusCode)
			}

			var wantData TestServerResponse
			wantBytes, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("%v", err)
			}
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
