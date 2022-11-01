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

func shouldbeMatchedRetried(t *testing.T, expectedRetried int, actualRetried int) error {
	if expectedRetried != actualRetried {
		return fmt.Errorf("expected retried: %d, but actual retried: %d", expectedRetried, actualRetried)
	}

	return nil
}

func shouldbeMatchedStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) error {
	if expectedStatusCode != actualStatusCode {
		return fmt.Errorf("expected status code: %d, but actual status code: %d", expectedStatusCode, actualStatusCode)
	}

	return nil
}

func ShouldbeMatchedData(t *testing.T, expectedData interface{}, actualBody io.Reader) error {
	var responseData TestServerResponse
	bodybytes, err := io.ReadAll(actualBody)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	err = json.Unmarshal(bodybytes, &responseData)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	var wantData TestServerResponse
	wantBytes, err := json.Marshal(expectedData)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = json.Unmarshal(wantBytes, &wantData)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if !reflect.DeepEqual(responseData, wantData) {
		responseDataBytes, _ := json.Marshal(responseData)
		wantBytes, _ := json.Marshal(expectedData)
		return fmt.Errorf("Retry.RetryRequest() = %v, want %v", string(responseDataBytes), string(wantBytes))
	}

	return nil
}

func TestRetry_RetryRequest_When_ServerHasSleep(t *testing.T) {
	type fields struct {
		RetryInterval int
		RetryMax      int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name              string
		fields            fields
		argsFn            func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args              *args
		serverSleepTimeMs int
		clientTimeoutMs   int
		want              *TestServerResponse
		wantStatusCode    int //statusc ode
		wantError         bool
		retried           int
	}{
		// Given Server's sleep set 500ms
		// And Server response code is 200
		// And Client's timeout not set
		// And data is not nil
		// Then OK
		// Then retried is 0
		{
			name: "1. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
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
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			retried:           0,
		},

		// Given Server's sleep set 500ms per a request
		// And Client's timeout 200ms
		// When data is not nil
		// Then error
		// Then retried is 3
		{
			name: "2. should be error",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
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
			wantStatusCode:    500,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			retried:           3,
			wantError:         true,
		},
		// Given Server's sleep set 1000ms
		// And Client's timeout 200ms
		// When data is not nil
		// Then error
		// Then retried is 3
		{
			name: "3. should be error",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
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
			wantStatusCode:    200,
			serverSleepTimeMs: 1000,
			clientTimeoutMs:   200,
			retried:           3,
			wantError:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var result *TestServerResponse
				if 0 < tt.serverSleepTimeMs {
					time.Sleep(time.Duration(tt.serverSleepTimeMs) * time.Millisecond)
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
				Policy: &RetryPolicy{
					RetryMax: tt.fields.RetryMax,
					Base:     tt.fields.RetryInterval,
					Cap:      tt.fields.RetryInterval,
				},
			}
			dataBytes, _ := json.Marshal(tt.args.data)
			serverUrl := server.URL

			request := tt.argsFn(tt.args.method, serverUrl, bytes.NewReader(dataBytes))
			got, err := r.Do(c, request, dataBytes)

			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}
			err = shouldbeMatchedRetried(t, tt.retried, r.retried)
			if err != nil {
				t.Error(err)
			}
			if tt.wantError {
				return
			}
			err = shouldbeMatchedStatusCode(t, tt.wantStatusCode, got.StatusCode)
			if err != nil {
				t.Error(err)
			}

			err = ShouldbeMatchedData(t, tt.want, got.Body)
			if err != nil {
				t.Error(err)
			}
			defer got.Body.Close()
		})
	}
}

func TestRetry_RetryRequest_When_ServerStatusCode500(t *testing.T) {
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
		triggerRetry    bool
		clientTimeoutMs int
		want            *TestServerResponse
		wantStatusCode  int //statusc ode
		wantError       bool
		retried         int
	}{
		// Given Server's 500 of status code per a request,
		// but at the last request, 200 of status code
		// And Client's timeout not set
		// When data is not nil
		// Then status code 200
		// Then retried is 3
		// Then should match req data and res data
		{
			name: "1. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
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
			retried:        3,
			triggerRetry:   true,
		},
		// Given Server's 500 of status code per a request,
		// but at the last request, 200 of status code
		// And Client's timeout not set
		// When data is nil
		// Then status code 200
		// Then retried is 3
		// Then should match req data and res data
		{
			name: "2. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode: 200,
			retried:        3,
			triggerRetry:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggerRetried := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var result *TestServerResponse
				triggerRetried += 1
				// When using Do, first a request is not retry.
				// It means at the first a request failed, then Do function should try retry
				if tt.triggerRetry && triggerRetried <= tt.fields.RetryMax || tt.wantError {
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
				Policy: &RetryPolicy{
					RetryMax: tt.fields.RetryMax,
					Base:     tt.fields.RetryInterval,
					Cap:      tt.fields.RetryInterval,
				},
			}
			dataBytes, _ := json.Marshal(tt.args.data)
			serverUrl := server.URL

			request := tt.argsFn(tt.args.method, serverUrl, bytes.NewReader(dataBytes))
			got, err := r.Do(c, request, dataBytes)

			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}
			err = shouldbeMatchedRetried(t, tt.retried, r.retried)
			if err != nil {
				t.Error(err)
			}
			if tt.wantError {
				return
			}

			err = shouldbeMatchedStatusCode(t, tt.wantStatusCode, got.StatusCode)
			if err != nil {
				t.Error(err)
			}
			err = ShouldbeMatchedData(t, tt.want, got.Body)
			if err != nil {
				t.Error(err)
			}
			defer got.Body.Close()
		})
	}
}

func TestRetry_RetryRequest_When_ServerHasSleep_But_LastRequestNoSleep(t *testing.T) {
	type fields struct {
		RetryInterval int
		RetryMax      int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name              string
		fields            fields
		argsFn            func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args              *args
		triggerRetry      bool
		serverSleepTimeMs int
		clientTimeoutMs   int
		want              *TestServerResponse
		respDataWhenRetry *TestServerResponse
		wantStatusCode    int //statusc ode
		wantError         bool
		retried           int
	}{
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And data is not nil
		// Then OK
		// Then retried is 3
		{
			name: "1. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data: &TestData{
					Name:    "None",
					Message: "NoneMessage",
				},
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			clientTimeoutMs:   200,
			serverSleepTimeMs: 500,
			retried:           3,
			triggerRetry:      true,
		},
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And data is nil
		// Then OK
		// Then retried is 3
		{
			name: "2. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data: &TestData{
					Name:    "None",
					Message: "NoneMessage",
				},
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			retried:           3,
			triggerRetry:      true,
		},
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout not set
		// And data is nil
		// Then OK
		// Then retried is 0
		{
			name: "3. should be ok",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			retried:           0,
			triggerRetry:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggerRetried := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var result *TestServerResponse
				triggerRetried += 1
				// When using Do, first a request is not retry.
				// It means at the first a request failed, then Do function should try retry
				if tt.triggerRetry && triggerRetried <= tt.fields.RetryMax || tt.wantError {
					time.Sleep(time.Duration(tt.serverSleepTimeMs) * time.Millisecond)
					w.WriteHeader(tt.wantStatusCode)

					result = tt.respDataWhenRetry
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
				Policy: &RetryPolicy{
					RetryMax: tt.fields.RetryMax,
					Base:     tt.fields.RetryInterval,
					Cap:      tt.fields.RetryInterval,
				},
			}
			dataBytes, _ := json.Marshal(tt.args.data)

			request := tt.argsFn(tt.args.method, server.URL, bytes.NewReader(dataBytes))
			got, err := r.Do(c, request, dataBytes)
			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}

			err = shouldbeMatchedRetried(t, tt.retried, r.retried)
			if err != nil {
				t.Error(err)
			}
			if tt.wantError {
				return
			}
			err = shouldbeMatchedStatusCode(t, tt.wantStatusCode, got.StatusCode)
			if err != nil {
				t.Error(err)
			}
			err = ShouldbeMatchedData(t, tt.want, got.Body)
			if err != nil {
				t.Error(err)
			}

			defer got.Body.Close()
		})
	}
}

func TestRetry_RetryRequest_When_UnknownUrl(t *testing.T) {
	type fields struct {
		RetryInterval int
		RetryMax      int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name              string
		fields            fields
		argsFn            func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args              *args
		triggerRetry      bool
		serverSleepTimeMs int
		clientTimeoutMs   int
		want              *TestServerResponse
		wantStatusCode    int //statusc ode
		wantError         bool
		manualServerUrl   string
		retried           int
	}{
		// Given unknown host
		// And Client's timeout 200ms per a request
		// And data nil
		// Then error
		// Then retried is 3
		{
			name: "1. should be error",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			want:              nil,
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			wantError:         true,
			retried:           3,
			manualServerUrl:   "http://127.0.0.10:9090",
		},

		// Given unknown host
		// And Client's timeout 200ms per a request
		// And data nil
		// Then error
		// Then retried is 3
		{
			name: "2. should be error",
			fields: fields{
				RetryInterval: 100,
				RetryMax:      3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			want:              nil,
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			wantError:         true,
			retried:           3,
			manualServerUrl:   "htpp://localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &http.Client{
				Timeout: time.Duration(tt.clientTimeoutMs) * time.Millisecond,
			}
			r := &Retry{
				Policy: &RetryPolicy{
					RetryMax: tt.fields.RetryMax,
					Base:     tt.fields.RetryInterval,
					Cap:      tt.fields.RetryInterval,
				},
			}
			dataBytes, _ := json.Marshal(tt.args.data)

			request := tt.argsFn(tt.args.method, tt.manualServerUrl, bytes.NewReader(dataBytes))
			_, err := r.Do(c, request, dataBytes)

			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}
			if tt.wantError {
				return
			}

			shouldbeMatchedRetried(t, tt.retried, r.retried)
			if err != nil {
				t.Error(err)
			}

		})
	}
}

func TestRetry_RetryRequest_Given_Jitter_When_ServerHasSleep_But_LastRequestNoSleep(t *testing.T) {
	type fields struct {
		Base       int
		Cap        int
		PolicyName RetryPolicyName
		RetryMax   int
	}

	type args struct {
		method string
		data   interface{}
	}
	tests := []struct {
		name              string
		fields            fields
		argsFn            func(string, string, io.Reader) *http.Request //method, url, io.Reader
		args              *args
		triggerRetry      bool
		serverSleepTimeMs int
		clientTimeoutMs   int
		want              *TestServerResponse
		respDataWhenRetry *TestServerResponse
		wantStatusCode    int //statusc ode
		wantError         bool
		retried           int
	}{
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And BackoffEqualJitter
		// And data is not nil
		// Then OK
		// Then retried is 3
		{
			name: "1. should be ok",
			fields: fields{
				PolicyName: RetryPolicyExpoEqualJitter,
				Base:       100,
				Cap:        300,
				RetryMax:   3,
			},
			args: &args{
				method: "GET",
				data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data: &TestData{
					Name:    "None",
					Message: "NoneMessage",
				},
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data: &TestData{
					Name:    "Hello",
					Message: "Message",
				},
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			clientTimeoutMs:   200,
			serverSleepTimeMs: 500,
			retried:           3,
			triggerRetry:      true,
		},
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And data is nil
		// And BackoffFullyJitter
		// Then OK
		// Then retried is 3
		{
			name: "2. should be ok",
			fields: fields{
				PolicyName: RetryPolicyExpoFullyJitter,
				Base:       100,
				Cap:        300,
				RetryMax:   3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data: &TestData{
					Name:    "None",
					Message: "NoneMessage",
				},
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			retried:           3,
			triggerRetry:      true,
		},
		// Given Server's sleep 500ms per a request, but it does not sleep at the a last request
		// And Client's timeout not set
		// And ExpoBackoff
		// And data is nil
		// Then OK
		// Then retried is 0
		{
			name: "3. should be ok",
			fields: fields{
				PolicyName: RetryPolicyExpoBackOff,
				Base:       100,
				Cap:        300,
				RetryMax:   3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			retried:           0,
			triggerRetry:      true,
		},
		// Given Server's sleep 300ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And ExpoBackoff
		// And data is nil
		// Then OK
		// Then retried is 3
		{
			name: "4. should be ok",
			fields: fields{
				PolicyName: RetryPolicyExpoBackOff,
				Base:       100,
				Cap:        200,
				RetryMax:   3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			retried:           3,
			triggerRetry:      true,
		},

		// Given Server's sleep 300ms per a request, but it does not sleep at the a last request
		// And Client's timeout 200ms
		// And NoExpoBackoff
		// And data is nil
		// Then OK
		// Then retried is 3
		{
			name: "5. should be ok",
			fields: fields{
				PolicyName: RetryPolicyNoBackOff,
				Base:       100,
				Cap:        100,
				RetryMax:   3,
			},
			args: &args{
				method: "GET",
				data:   nil,
			},
			argsFn: func(method, url string, r io.Reader) *http.Request {
				req, _ := http.NewRequest(method, url, r)

				return req
			},
			respDataWhenRetry: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			want: &TestServerResponse{
				Data:         nil,
				ErrorMessage: nil,
			},
			wantStatusCode:    200,
			serverSleepTimeMs: 500,
			clientTimeoutMs:   200,
			retried:           3,
			triggerRetry:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggerRetried := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var result *TestServerResponse
				triggerRetried += 1
				// When using Do, first a request is not retry.
				// It means at the first a request failed, then Do function should try retry
				if tt.triggerRetry && triggerRetried <= tt.fields.RetryMax || tt.wantError {
					time.Sleep(time.Duration(tt.serverSleepTimeMs) * time.Millisecond)
					w.WriteHeader(tt.wantStatusCode)

					result = tt.respDataWhenRetry
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
				Policy: &RetryPolicy{
					PolicyName: tt.fields.PolicyName,
					RetryMax:   tt.fields.RetryMax,
					Base:       tt.fields.Base,
					Cap:        tt.fields.Cap,
				},
			}
			dataBytes, _ := json.Marshal(tt.args.data)

			request := tt.argsFn(tt.args.method, server.URL, bytes.NewReader(dataBytes))
			got, err := r.Do(c, request, dataBytes)
			if (err != nil) != tt.wantError {
				t.Errorf("%v", err)
				return
			}

			err = shouldbeMatchedRetried(t, tt.retried, r.retried)
			if err != nil {
				t.Error(err)
			}
			if tt.wantError {
				return
			}
			err = shouldbeMatchedStatusCode(t, tt.wantStatusCode, got.StatusCode)
			if err != nil {
				t.Error(err)
			}
			err = ShouldbeMatchedData(t, tt.want, got.Body)
			if err != nil {
				t.Error(err)
			}

			defer got.Body.Close()
		})
	}
}
