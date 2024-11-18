package task

import (
	"net"
	"testing"
)

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "valid ipv4",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "valid ipv6",
			ip:   "2001:db8::1",
			want: false,
		},
		{
			name: "empty string",
			ip:   "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIPv4(tt.ip); got != tt.want {
				t.Errorf("isIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPRanges_FixIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantIP  string
		wantv4  bool
	}{
		{
			name:    "ipv4 without mask",
			ip:      "192.168.1.1",
			wantIP:  "192.168.1.1/32",
			wantv4:  true,
		},
		{
			name:    "ipv4 with mask",
			ip:      "192.168.1.0/24",
			wantIP:  "192.168.1.0/24",
			wantv4:  true,
		},
		{
			name:    "ipv6 without mask",
			ip:      "2001:db8::1",
			wantIP:  "2001:db8::1/128",
			wantv4:  false,
		},
		{
			name:    "ipv6 with mask",
			ip:      "2001:db8::/64",
			wantIP:  "2001:db8::/64",
			wantv4:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newIPRanges()
			got := r.fixIP(tt.ip)
			if got != tt.wantIP {
				t.Errorf("fixIP() = %v, want %v", got, tt.wantIP)
			}
		})
	}
}

func TestIPRanges_GetIPRange(t *testing.T) {
	tests := []struct {
		name     string
		setupIP  string
		wantMin  byte
		wantHost byte
	}{
		{
			name:     "ipv4 /32",
			setupIP:  "192.168.1.1/32",
			wantMin:  1,
			wantHost: 0,
		},
		{
			name:     "ipv4 /24",
			setupIP:  "192.168.1.0/24",
			wantMin:  0,
			wantHost: 255,
		},
		{
			name:     "ipv4 /16",
			setupIP:  "192.168.0.0/16",
			wantMin:  0,
			wantHost: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newIPRanges()
			r.parseCIDR(r.fixIP(tt.setupIP))
			gotMin, gotHosts := r.getIPRange()
			if gotMin != tt.wantMin || gotHosts != tt.wantHost {
				t.Errorf("getIPRange() = (%v, %v), want (%v, %v)",
					gotMin, gotHosts, tt.wantMin, tt.wantHost)
			}
		})
	}
}

func TestIPRanges_AppendIP(t *testing.T) {
	r := newIPRanges()
	testIP := net.ParseIP("192.168.1.1")
	
	r.appendIP(testIP)
	
	if len(r.ips) != 1 {
		t.Errorf("appendIP() did not append IP, got len = %v, want 1", len(r.ips))
	}
	
	if !r.ips[0].IP.Equal(testIP) {
		t.Errorf("appendIP() appended wrong IP, got %v, want %v", r.ips[0].IP, testIP)
	}
}
