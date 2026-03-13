package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateRobot_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2.0/robots" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			t.Errorf("unexpected auth: user=%q pass=%q ok=%v", user, pass, ok)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		var req robotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Name != "my-cluster" {
			t.Errorf("robot name = %q, want %q", req.Name, "my-cluster")
		}
		if req.Level != "system" {
			t.Errorf("robot level = %q, want %q", req.Level, "system")
		}
		if req.Duration != -1 {
			t.Errorf("robot duration = %d, want %d", req.Duration, -1)
		}
		if len(req.Permissions) != 2 {
			t.Errorf("permissions count = %d, want 2", len(req.Permissions))
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(robotResponse{
			ID:     42,
			Name:   "robot$my-cluster",
			Secret: "robot-secret-token",
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	c.password = "secret" // skip K8s secret read

	creds, err := c.createRobot(t.Context(), "my-cluster", "secret")
	if err != nil {
		t.Fatalf("createRobot: %v", err)
	}
	if creds.Name != "robot$my-cluster" {
		t.Errorf("creds.Name = %q, want %q", creds.Name, "robot$my-cluster")
	}
	if creds.Secret != "robot-secret-token" {
		t.Errorf("creds.Secret = %q, want %q", creds.Secret, "robot-secret-token")
	}
}

func TestCreateRobot_Conflict_DeleteAndRetry(t *testing.T) {
	var (
		createCalls  int
		deleteCalled bool
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2.0/robots":
			createCalls++
			if createCalls == 1 {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(robotResponse{
				ID:     99,
				Name:   "robot$test-cluster",
				Secret: "new-secret",
			})

		case r.Method == http.MethodGet && r.URL.Path == "/api/v2.0/robots":
			q := r.URL.Query().Get("q")
			if !strings.Contains(q, "test-cluster") {
				t.Errorf("unexpected query: %q", q)
			}
			json.NewEncoder(w).Encode([]robotListItem{
				{ID: 50, Name: "robot$test-cluster"},
			})

		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2.0/robots/50":
			deleteCalled = true
			w.WriteHeader(http.StatusOK)

		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	c.password = "secret"

	creds, err := c.EnsureRobotAccount(t.Context(), "test-cluster")
	if err != nil {
		t.Fatalf("EnsureRobotAccount: %v", err)
	}

	if !deleteCalled {
		t.Error("expected DELETE to be called")
	}
	if createCalls != 2 {
		t.Errorf("createCalls = %d, want 2", createCalls)
	}
	if creds.Name != "robot$test-cluster" {
		t.Errorf("creds.Name = %q, want %q", creds.Name, "robot$test-cluster")
	}
	if creds.Secret != "new-secret" {
		t.Errorf("creds.Secret = %q, want %q", creds.Secret, "new-secret")
	}
}

func TestCreateRobot_UnexpectedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	c.password = "secret"

	_, err := c.EnsureRobotAccount(context.Background(), "fail-cluster")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status 500") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFindRobotID_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]robotListItem{})
	}))
	defer srv.Close()

	c := NewClient(srv.URL)

	_, err := c.findRobotID(t.Context(), "missing-cluster", "secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRobotPermissions(t *testing.T) {
	var gotReq robotRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(robotResponse{ID: 1, Name: "robot$perm-test", Secret: "s"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	c.password = "secret"

	_, err := c.EnsureRobotAccount(t.Context(), "perm-test")
	if err != nil {
		t.Fatalf("EnsureRobotAccount: %v", err)
	}

	if len(gotReq.Permissions) != 2 {
		t.Fatalf("permissions count = %d, want 2", len(gotReq.Permissions))
	}

	systemPerm := gotReq.Permissions[0]
	if systemPerm.Kind != "system" || systemPerm.Namespace != "/" {
		t.Errorf("system permission: kind=%q namespace=%q", systemPerm.Kind, systemPerm.Namespace)
	}
	wantSystemActions := map[string]bool{"create": true, "list": true, "read": true}
	for _, a := range systemPerm.Access {
		if a.Resource != "project" {
			t.Errorf("system access resource = %q, want %q", a.Resource, "project")
		}
		delete(wantSystemActions, a.Action)
	}
	if len(wantSystemActions) > 0 {
		t.Errorf("missing system project actions: %v", wantSystemActions)
	}

	projectPerm := gotReq.Permissions[1]
	if projectPerm.Kind != "project" || projectPerm.Namespace != "*" {
		t.Errorf("project permission: kind=%q namespace=%q", projectPerm.Kind, projectPerm.Namespace)
	}

	wantProjectAccess := map[string][]string{
		"member": {"create", "update", "list", "read", "delete"},
		"robot":  {"create", "list", "read"},
	}
	gotAccess := map[string]map[string]bool{}
	for _, a := range projectPerm.Access {
		if gotAccess[a.Resource] == nil {
			gotAccess[a.Resource] = map[string]bool{}
		}
		gotAccess[a.Resource][a.Action] = true
	}
	for resource, actions := range wantProjectAccess {
		for _, action := range actions {
			if !gotAccess[resource][action] {
				t.Errorf("missing project access: resource=%q action=%q", resource, action)
			}
		}
	}
}

func TestDeleteRobot_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v2.0/robots/42" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)

	err := c.deleteRobot(t.Context(), 42, "secret")
	if err != nil {
		t.Fatalf("deleteRobot: %v", err)
	}
}

func TestDeleteRobot_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)

	err := c.deleteRobot(t.Context(), 42, "secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("unexpected status %d", http.StatusForbidden)) {
		t.Errorf("unexpected error: %v", err)
	}
}
