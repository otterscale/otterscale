package core

import (
	"strings"
	"testing"
)

func TestParseRancherProjectID(t *testing.T) {
	validCluster := strings.Repeat("a", 63)
	validProject := strings.Repeat("b", 253)
	tests := []struct {
		name        string
		id          string
		wantCluster string
		wantProject string
		wantErr     bool
	}{
		{name: "rancher cluster", id: "c-m-abcde:p-vwxyz", wantCluster: "c-m-abcde", wantProject: "p-vwxyz"},
		{name: "local cluster", id: "local:p-local1", wantCluster: "local", wantProject: "p-local1"},
		{name: "project subdomain", id: "local:team.project", wantCluster: "local", wantProject: "team.project"},
		{name: "segment boundaries", id: validCluster + ":" + validProject, wantCluster: validCluster, wantProject: validProject},
		{name: "empty", id: "", wantErr: true},
		{name: "missing cluster", id: ":p-test", wantErr: true},
		{name: "missing project", id: "local:", wantErr: true},
		{name: "multiple separators", id: "local:p-test:extra", wantErr: true},
		{name: "cluster subdomain", id: "cluster.example:p-test", wantErr: true},
		{name: "uppercase", id: "local:P-test", wantErr: true},
		{name: "cluster too long", id: strings.Repeat("a", 64) + ":p-test", wantErr: true},
		{name: "project too long", id: "local:" + strings.Repeat("b", 254), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cluster, project, err := ParseRancherProjectID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRancherProjectID: %v", err)
			}
			if cluster != tt.wantCluster || project != tt.wantProject {
				t.Fatalf("got %q/%q, want %q/%q", cluster, project, tt.wantCluster, tt.wantProject)
			}
		})
	}
}
