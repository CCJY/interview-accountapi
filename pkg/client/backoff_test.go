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
				base:       tt.base, cap: tt.maxDuration}
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
				base:       tt.base, cap: tt.maxDuration}
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
				base:       tt.base, cap: tt.maxDuration}
			sleeptime := backoff.CalcuateSleep(tt.retried, 0)
			if tt.maxDuration < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}
