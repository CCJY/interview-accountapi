package client

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRetryWithBackoff(t *testing.T) {
	tests := []struct {
		base        int
		maxDuration int
		retried     int
	}{
		{
			base:        100,
			maxDuration: 30000,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 50000,
			retried:     3,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     2,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     0,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			rand.Seed(time.Now().UnixNano())
			backoff := RetryPolicy{
				PolicyName: RetryPolicyExpoBackOff,
				Base:       tt.base, Cap: tt.maxDuration}
			sleeptime := backoff.CalcuateSleep(tt.retried, 0)

			if tt.maxDuration < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}

func TestRetryWithEqualJitter(t *testing.T) {
	tests := []struct {
		base        int
		maxDuration int
		retried     int
	}{
		{
			base:        100,
			maxDuration: 50000,
			retried:     3,
		},
		{
			base:        100,
			maxDuration: 30000,
			retried:     2,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     3,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			rand.Seed(time.Now().UnixNano())
			backoff := RetryPolicy{
				PolicyName: RetryPolicyExpoEqualJitter,
				Base:       tt.base, Cap: tt.maxDuration}
			sleeptime := backoff.CalcuateSleep(tt.retried, 0)

			if tt.maxDuration < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}

func TestRetryWithFullyJitter(t *testing.T) {
	tests := []struct {
		base        int
		maxDuration int
		retried     int
	}{
		{
			base:        100,
			maxDuration: 1000,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 1000,
			retried:     2,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     3,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			backoff := RetryPolicy{
				PolicyName: RetryPolicyExpoFullyJitter,
				Base:       tt.base, Cap: tt.maxDuration}
			sleeptime := backoff.CalcuateSleep(tt.retried, 0)
			if tt.maxDuration < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}

func TestRetryWithDecorrJitter(t *testing.T) {
	tests := []struct {
		base        int
		maxDuration int
		retried     int
	}{
		{
			base:        100,
			maxDuration: 100,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 1000,
			retried:     2,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     3,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			backoff := RetryPolicy{
				PolicyName: RetryPolicyExpoDecorrJitter,
				Base:       tt.base, Cap: tt.maxDuration}
			sleeptime := backoff.CalcuateSleep(tt.retried, backoff.Base)
			if tt.maxDuration < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}
