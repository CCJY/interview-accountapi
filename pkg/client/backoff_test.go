package client

import (
	"fmt"
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
			maxDuration: 300,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 300,
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
			sleeptime := Backoff(tt.base, tt.maxDuration, tt.retried)

			if sleeptime < tt.base || tt.maxDuration < sleeptime {
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
			maxDuration: 300,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     1,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			sleeptime := EqualJitter(tt.base, tt.maxDuration, tt.retried)

			if sleeptime < tt.base || tt.maxDuration < sleeptime {
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
			maxDuration: 300,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     1,
		},
		{
			base:        100,
			maxDuration: 300,
			retried:     2,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			sleeptime := FullyJitter(tt.base, tt.maxDuration, tt.retried)
			if tt.base < sleeptime {
				t.Errorf("wrong: %v", sleeptime)
			}
			fmt.Println(time.Duration(sleeptime) * time.Millisecond)
		})
	}
}
