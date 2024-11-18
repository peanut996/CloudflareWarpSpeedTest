package utils

import (
	"net"
	"testing"
	"time"
)

func TestCloudflareIPData_getLossRate(t *testing.T) {
	tests := []struct {
		name     string
		pingData *PingData
		want     float32
	}{
		{
			name: "no loss",
			pingData: &PingData{
				Sent:     10,
				Received: 10,
			},
			want: 0,
		},
		{
			name: "50% loss",
			pingData: &PingData{
				Sent:     10,
				Received: 5,
			},
			want: 0.5,
		},
		{
			name: "100% loss",
			pingData: &PingData{
				Sent:     10,
				Received: 0,
			},
			want: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := &CloudflareIPData{
				PingData: tt.pingData,
			}
			if got := cf.getLossRate(); got != tt.want {
				t.Errorf("CloudflareIPData.getLossRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPingDelaySet_FilterDelay(t *testing.T) {
	// Save original values
	origMaxDelay := InputMaxDelay
	origMinDelay := InputMinDelay
	defer func() {
		InputMaxDelay = origMaxDelay
		InputMinDelay = origMinDelay
	}()

	// Set test values
	InputMaxDelay = 100 * time.Millisecond
	InputMinDelay = 10 * time.Millisecond

	testIP, _ := net.ResolveUDPAddr("udp", "1.1.1.1:0")
	
	tests := []struct {
		name string
		set  PingDelaySet
		want int // number of items that should remain after filtering
	}{
		{
			name: "all within range",
			set: PingDelaySet{
				{PingData: &PingData{IP: testIP, Delay: 50 * time.Millisecond}},
				{PingData: &PingData{IP: testIP, Delay: 75 * time.Millisecond}},
			},
			want: 2,
		},
		{
			name: "some outside range",
			set: PingDelaySet{
				{PingData: &PingData{IP: testIP, Delay: 5 * time.Millisecond}},
				{PingData: &PingData{IP: testIP, Delay: 50 * time.Millisecond}},
				{PingData: &PingData{IP: testIP, Delay: 150 * time.Millisecond}},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.FilterDelay(); len(got) != tt.want {
				t.Errorf("PingDelaySet.FilterDelay() returned %v items, want %v", len(got), tt.want)
			}
		})
	}
}

func TestPingDelaySet_FilterLossRate(t *testing.T) {
	// Save original value
	origMaxLossRate := InputMaxLossRate
	defer func() {
		InputMaxLossRate = origMaxLossRate
	}()

	// Set test value
	InputMaxLossRate = 0.5

	testIP, _ := net.ResolveUDPAddr("udp", "1.1.1.1:0")
	
	tests := []struct {
		name string
		set  PingDelaySet
		want int // number of items that should remain after filtering
	}{
		{
			name: "all within range",
			set: PingDelaySet{
				{PingData: &PingData{IP: testIP, Sent: 10, Received: 8}}, // 20% loss
				{PingData: &PingData{IP: testIP, Sent: 10, Received: 7}}, // 30% loss
			},
			want: 2,
		},
		{
			name: "some outside range",
			set: PingDelaySet{
				{PingData: &PingData{IP: testIP, Sent: 10, Received: 8}}, // 20% loss
				{PingData: &PingData{IP: testIP, Sent: 10, Received: 4}}, // 60% loss
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.FilterLossRate(); len(got) != tt.want {
				t.Errorf("PingDelaySet.FilterLossRate() returned %v items, want %v", len(got), tt.want)
			}
		})
	}
}
