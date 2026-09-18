package harbor

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	testCluster  = "prod"
	testSecret   = "OtderivedSecret0"
	testPassword = "admin-password"
)

// request records one call, so a test can assert on the sequence.
type request struct {
	method string
	path   string
	query  string
	user   string
	pass   string
	body   []byte
}

// recorder logs every request and answers from a per-route table.
type recorder struct {
	t        *testing.T
	requests []request
	// handlers is keyed by "METHOD /path"; a missing route fails the test.
	handlers map[string]http.HandlerFunc
}

func (r *recorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.t.Helper()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		r.t.Fatalf("read request body: %v", err)
	}
	user, pass, _ := req.BasicAuth()
	r.requests = append(r.requests, request{
		method: req.Method,
		path:   req.URL.Path,
		query:  req.URL.RawQuery,
		user:   user,
		pass:   pass,
		body:   body,
	})

	key := req.Method + " " + req.URL.Path
	handler, ok := r.handlers[key]
	if !ok {
		r.t.Errorf("unexpected request %s", key)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler(w, req)
}

func (r *recorder) calls() []string {
	out := make([]string, 0, len(r.requests))
	for _, req := range r.requests {
		out = append(out, req.method+" "+req.path)
	}
	return out
}

func newTestClient(t *testing.T, handlers map[string]http.HandlerFunc) (*Client, *recorder) {
	t.Helper()

	rec := &recorder{t: t, handlers: handlers}
	server := httptest.NewServer(rec)
	t.Cleanup(server.Close)

	return NewClient(server.URL, func() (string, error) { return testPassword, nil }), rec
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func TestClient_EnsureRobotAccount_Creates(t *testing.T) {
	client, rec := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			// Harbor echoes the prefixed name; it may not echo the secret, which
			// is why the client trusts the one it sent.
			writeJSON(t, w, http.StatusCreated, robotResponse{ID: 7, Name: "robot$" + testCluster})
		},
	})

	creds, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
	if err != nil {
		t.Fatalf("EnsureRobotAccount() error = %v", err)
	}
	if creds.Name != "robot$"+testCluster {
		t.Errorf("name = %q, want %q", creds.Name, "robot$"+testCluster)
	}
	if creds.Secret != testSecret {
		t.Errorf("secret = %q, want the one we supplied %q", creds.Secret, testSecret)
	}

	if len(rec.requests) != 1 {
		t.Fatalf("calls = %v, want one create", rec.calls())
	}

	var sent robotRequest
	if err := json.Unmarshal(rec.requests[0].body, &sent); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sent.Name != testCluster {
		t.Errorf("name = %q, want the bare cluster name %q", sent.Name, testCluster)
	}
	if sent.Secret != testSecret {
		t.Errorf("secret = %q, want %q", sent.Secret, testSecret)
	}
	if sent.Duration != robotNeverExpires {
		t.Errorf("duration = %d, want %d", sent.Duration, robotNeverExpires)
	}
	if sent.Level != robotLevel {
		t.Errorf("level = %q, want %q", sent.Level, robotLevel)
	}
	if len(sent.Permissions) != 2 {
		t.Fatalf("permissions = %d scopes, want 2", len(sent.Permissions))
	}
	if sent.Permissions[0].Namespace != "/" || sent.Permissions[1].Namespace != "*" {
		t.Errorf("permission scopes = %q/%q, want //*",
			sent.Permissions[0].Namespace, sent.Permissions[1].Namespace)
	}
}

// The behavioral change from the implementation this replaces, which deleted
// and re-created the account and so revoked the credential in use on every
// call.
func TestClient_EnsureRobotAccount_AdoptsExisting(t *testing.T) {
	client, rec := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusConflict)
		},
		"GET " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusOK, []robotListItem{
				{ID: 3, Name: "robot$eu-" + testCluster},
				{ID: 9, Name: "robot$" + testCluster},
			})
		},
		"PATCH " + robotsPath + "/9": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	})

	creds, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
	if err != nil {
		t.Fatalf("EnsureRobotAccount() error = %v", err)
	}
	if creds.Secret != testSecret {
		t.Errorf("secret = %q, want %q", creds.Secret, testSecret)
	}

	want := []string{
		"POST " + robotsPath,
		"GET " + robotsPath,
		"PATCH " + robotsPath + "/9",
	}
	got := rec.calls()
	if len(got) != len(want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("calls = %v, want %v", got, want)
		}
	}
	for _, call := range got {
		if strings.HasPrefix(call, http.MethodDelete) {
			t.Error("the existing robot must never be deleted")
		}
	}

	var sent robotSecretRequest
	if err := json.Unmarshal(rec.requests[2].body, &sent); err != nil {
		t.Fatalf("unmarshal patch body: %v", err)
	}
	if sent.Secret != testSecret {
		t.Errorf("patched secret = %q, want %q", sent.Secret, testSecret)
	}
}

// The list query is a substring match, so "prod" must not adopt "eu-prod".
func TestClient_EnsureRobotAccount_MatchesExactName(t *testing.T) {
	client, rec := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusConflict)
		},
		"GET " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusOK, []robotListItem{
				{ID: 3, Name: "robot$eu-" + testCluster},
			})
		},
	})

	_, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), testCluster) {
		t.Errorf("error = %v, want it to name the cluster", err)
	}
	for _, call := range rec.calls() {
		if strings.HasPrefix(call, http.MethodPatch) {
			t.Error("a robot belonging to another cluster must not be patched")
		}
	}
}

func TestClient_EnsureRobotAccount_SendsBasicAuth(t *testing.T) {
	client, rec := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusCreated, robotResponse{ID: 1, Name: "robot$" + testCluster})
		},
	})

	if _, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret); err != nil {
		t.Fatalf("EnsureRobotAccount() error = %v", err)
	}
	for _, req := range rec.requests {
		if req.user != adminUser || req.pass != testPassword {
			t.Errorf("%s %s auth = %q/%q, want %q/%q",
				req.method, req.path, req.user, req.pass, adminUser, testPassword)
		}
	}
}

// Regression test for the cache the previous implementation kept forever,
// which made rotating the admin password require a restart.
func TestClient_EnsureRobotAccount_ResolvesPasswordPerCall(t *testing.T) {
	rec := &recorder{t: t, handlers: map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusCreated, robotResponse{ID: 1, Name: "robot$" + testCluster})
		},
	}}
	server := httptest.NewServer(rec)
	t.Cleanup(server.Close)

	reads := 0
	passwords := []string{"first", "second"}
	client := NewClient(server.URL, func() (string, error) {
		password := passwords[min(reads, len(passwords)-1)]
		reads++
		return password, nil
	})

	for range 2 {
		if _, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret); err != nil {
			t.Fatalf("EnsureRobotAccount() error = %v", err)
		}
	}

	if reads != 2 {
		t.Errorf("password reads = %d, want one per call", reads)
	}
	if rec.requests[0].pass == rec.requests[1].pass {
		t.Error("the second call reused the first call's password")
	}
}

func TestClient_EnsureRobotAccount_MissingPassword(t *testing.T) {
	tests := []struct {
		name     string
		password func() (string, error)
	}{
		{
			name: "unreadable",
			password: func() (string, error) {
				return "", errors.New("open /etc/otterscale/harbor/admin-password: no such file")
			},
		},
		{
			name:     "empty",
			password: func() (string, error) { return "", nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := &recorder{t: t, handlers: map[string]http.HandlerFunc{}}
			server := httptest.NewServer(rec)
			t.Cleanup(server.Close)

			client := NewClient(server.URL, tt.password)

			_, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if code, ok := core.DomainErrorCode(err); !ok || code != core.ErrorCodeFailedPrecondition {
				t.Errorf("code = %v (domain=%t), want %v", code, ok, core.ErrorCodeFailedPrecondition)
			}
			if len(rec.requests) != 0 {
				t.Error("no request may be sent without a credential")
			}
		})
	}
}

func TestClient_EnsureRobotAccount_RejectedCredential(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			client, _ := newTestClient(t, map[string]http.HandlerFunc{
				"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, "unauthorized", status)
				},
			})

			_, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if code, ok := core.DomainErrorCode(err); !ok || code != core.ErrorCodeFailedPrecondition {
				t.Errorf("code = %v (domain=%t), want %v", code, ok, core.ErrorCodeFailedPrecondition)
			}
		})
	}
}

func TestClient_EnsureRobotAccount_SurfacesServerError(t *testing.T) {
	client, _ := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "the registry is on fire", http.StatusInternalServerError)
		},
	})

	_, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "the registry is on fire") {
		t.Errorf("error = %v, want it to carry Harbor's response body", err)
	}
}

func TestClient_EnsureRobotAccount_QueriesByName(t *testing.T) {
	client, rec := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusConflict)
		},
		"GET " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusOK, []robotListItem{{ID: 1, Name: "robot$" + testCluster}})
		},
		"PATCH " + robotsPath + "/1": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	})

	if _, err := client.EnsureRobotAccount(t.Context(), testCluster, testSecret); err != nil {
		t.Fatalf("EnsureRobotAccount() error = %v", err)
	}

	want := "q=" + "name%3D~" + testCluster
	if got := rec.requests[1].query; got != want {
		t.Errorf("list query = %q, want %q", got, want)
	}
}

func TestClient_EnsureRobotAccount_CancelledContext(t *testing.T) {
	client, _ := newTestClient(t, map[string]http.HandlerFunc{
		"POST " + robotsPath: func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(t, w, http.StatusCreated, robotResponse{ID: 1, Name: "robot$" + testCluster})
		},
	})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, err := client.EnsureRobotAccount(ctx, testCluster, testSecret); err == nil {
		t.Error("expected an error for a canceled context, got nil")
	}
}

func TestNewClient_TrimsTrailingSlash(t *testing.T) {
	client := NewClient("https://harbor.example.com:8443/", func() (string, error) { return testPassword, nil })

	if want := "https://harbor.example.com:8443"; client.baseURL != want {
		t.Errorf("baseURL = %q, want %q", client.baseURL, want)
	}
}

func TestHarborRobotName(t *testing.T) {
	if got, want := core.HarborRobotName("prod"), "robot$prod"; got != want {
		t.Errorf("HarborRobotName() = %q, want %q", got, want)
	}
}
