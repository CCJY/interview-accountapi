package client

import (
	"testing"
	"time"
)

func TestRetry_RetryRequest(t *testing.T) {
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
			got, got1 := r.RetryRequest()
			if got != tt.want {
				t.Errorf("Retry.RetryRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Retry.RetryRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
