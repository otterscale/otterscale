package server

import (
	"strings"
	"testing"
)

func TestResolveTunnelHost(t *testing.T) {
	tests := []struct {
		name              string
		externalTunnelURL string
		tunnelAddress     string
		want              string
		wantErr           string
	}{
		{
			name:              "external URL wins over listen address",
			externalTunnelURL: "https://tunnel.example.com:8300",
			tunnelAddress:     "0.0.0.0:8300",
			want:              "tunnel.example.com",
		},
		{
			name:              "external URL without port",
			externalTunnelURL: "https://tunnel.example.com",
			tunnelAddress:     "0.0.0.0:8300",
			want:              "tunnel.example.com",
		},
		{
			name:              "external URL with IP literal",
			externalTunnelURL: "https://203.0.113.10:8300",
			tunnelAddress:     "0.0.0.0:8300",
			want:              "203.0.113.10",
		},
		{
			name:              "external URL with no host is rejected",
			externalTunnelURL: "tunnel.example.com:8300",
			tunnelAddress:     "127.0.0.1:8300",
			wantErr:           "has no host",
		},
		{
			name:          "concrete listen address is used as fallback",
			tunnelAddress: "127.0.0.1:8300",
			want:          "127.0.0.1",
		},
		{
			name:          "IPv4 wildcard listen address is rejected",
			tunnelAddress: "0.0.0.0:8300",
			wantErr:       "wildcard address",
		},
		{
			name:          "IPv6 wildcard listen address is rejected",
			tunnelAddress: "[::]:8300",
			wantErr:       "wildcard address",
		},
		{
			name:          "empty host in listen address is rejected",
			tunnelAddress: ":8300",
			wantErr:       "wildcard address",
		},
		{
			name:          "unparsable listen address is rejected",
			tunnelAddress: "not-an-address",
			wantErr:       "parse tunnel address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				TunnelAddress:     tt.tunnelAddress,
				ExternalTunnelURL: tt.externalTunnelURL,
			}

			got, err := resolveTunnelHost(cfg)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got host %q", tt.wantErr, got)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("host = %q, want %q", got, tt.want)
			}
		})
	}
}
