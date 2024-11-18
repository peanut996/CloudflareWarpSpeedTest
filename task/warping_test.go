package task

import (
	"net"
	"testing"
	"time"

	"github.com/peanut996/CloudflareWarpSpeedTest/utils"
)

func TestUDPAddr_FullAddress(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		port int
		want string
	}{
		{
			name: "ipv4 address",
			ip:   "192.168.1.1",
			port: 8080,
			want: "192.168.1.1:8080",
		},
		{
			name: "ipv6 address",
			ip:   "2001:db8::1",
			port: 8080,
			want: "[2001:db8::1]:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipAddr := &net.IPAddr{IP: net.ParseIP(tt.ip)}
			addr := &UDPAddr{
				IP:   ipAddr,
				Port: tt.port,
			}
			if got := addr.FullAddress(); got != tt.want {
				t.Errorf("UDPAddr.FullAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUDPAddr_ToUDPAddr(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		port    int
		wantErr bool
	}{
		{
			name:    "valid ipv4",
			ip:      "192.168.1.1",
			port:    8080,
			wantErr: false,
		},
		{
			name:    "valid ipv6",
			ip:      "2001:db8::1",
			port:    8080,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipAddr := &net.IPAddr{IP: net.ParseIP(tt.ip)}
			addr := &UDPAddr{
				IP:   ipAddr,
				Port: tt.port,
			}
			got := addr.ToUDPAddr()
			if (got == nil) != tt.wantErr {
				t.Errorf("UDPAddr.ToUDPAddr() error = %v, wantErr %v", got == nil, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Port != tt.port {
				t.Errorf("UDPAddr.ToUDPAddr() port = %v, want %v", got.Port, tt.port)
			}
		})
	}
}

func TestWarping_AppendIPData(t *testing.T) {
	w := NewWarping()
	testIP := &net.IPAddr{IP: net.ParseIP("192.168.1.1")}
	pingData := &utils.PingData{
		IP:       &net.UDPAddr{IP: testIP.IP, Port: 8080},
		Sent:     10,
		Received: 8,
		Delay:    100 * time.Millisecond,
	}

	w.appendIPData(pingData)

	if len(w.csv) != 1 {
		t.Errorf("Warping.appendIPData() did not append data, got len = %v, want 1", len(w.csv))
		return
	}

	got := w.csv[0]
	if got.Delay != pingData.Delay {
		t.Errorf("Warping.appendIPData() delay = %v, want %v", got.Delay, pingData.Delay)
	}
	if got.Sent != pingData.Sent {
		t.Errorf("Warping.appendIPData() sent = %v, want %v", got.Sent, pingData.Sent)
	}
	if got.Received != pingData.Received {
		t.Errorf("Warping.appendIPData() received = %v, want %v", got.Received, pingData.Received)
	}
}

func TestEncodeBase64ToHex(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid base64",
			key:     warpPublicKey,  // Use the actual WARP public key which is known to be valid
			want:    "6e65ce0be17517110c17d77288ad87e7fd5252dcc7d09b95a39d61db03df832a",
			wantErr: false,
		},
		{
			name:    "invalid base64",
			key:     "invalid@@",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			key:     "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeBase64ToHex(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeBase64ToHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("encodeBase64ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNoisePublicKeyFromBase64(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid warp public key",
			key:     warpPublicKey,
			wantErr: false,
		},
		{
			name:    "invalid key",
			key:     "invalid@@",
			wantErr: true,
		},
		{
			name:    "empty string",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNoisePublicKeyFromBase64(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNoisePublicKeyFromBase64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
