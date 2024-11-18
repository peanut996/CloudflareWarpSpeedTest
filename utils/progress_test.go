package utils

import (
	"testing"
	"time"
)

func TestNewBar(t *testing.T) {
	tests := []struct {
		name      string
		count     int
		strStart  string
		strEnd    string
		wantCount int
	}{
		{
			name:      "basic progress bar",
			count:     100,
			strStart:  "Processing",
			strEnd:    "items",
			wantCount: 100,
		},
		{
			name:      "zero count",
			count:     0,
			strStart:  "Starting",
			strEnd:    "tasks",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBar(tt.count, tt.strStart, tt.strEnd)
			if bar == nil {
				t.Error("NewBar() returned nil")
			}
			if bar.pb == nil {
				t.Error("NewBar() returned bar with nil pb")
			}
		})
	}
}

func TestBar_Grow(t *testing.T) {
	tests := []struct {
		name        string
		initialSize int
		growBy     int
		strVal     string
	}{
		{
			name:        "grow by one",
			initialSize: 10,
			growBy:     1,
			strVal:     "processing item 1",
		},
		{
			name:        "grow by multiple",
			initialSize: 100,
			growBy:     5,
			strVal:     "batch processing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBar(tt.initialSize, "Start", "End")
			bar.Grow(tt.growBy, tt.strVal)
			// Allow a short time for the progress bar to update
			time.Sleep(10 * time.Millisecond)
		})
	}
}

func TestBar_Done(t *testing.T) {
	tests := []struct {
		name        string
		initialSize int
	}{
		{
			name:        "complete small bar",
			initialSize: 10,
		},
		{
			name:        "complete large bar",
			initialSize: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBar(tt.initialSize, "Start", "End")
			bar.Done()
			// Allow a short time for the progress bar to finish
			time.Sleep(10 * time.Millisecond)
		})
	}
}
