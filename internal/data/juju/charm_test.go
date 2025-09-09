package juju

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/wrap"
)

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
 * -------------------------------------------------------------------------- */
func charm_mustNewJuju(t testing.TB, cfg *config.Config) *Juju {
	j := New(cfg)

	return j
}

/* -------------------------------------------------------------------------- *
 * Minimal JSON payloads used by the HTTP‑mock server
 * -------------------------------------------------------------------------- */

/*
The real “find” endpoint returns a slice of full core.Charm objects.

	We only need the fields that the implementation accesses, so we give
	each charm a name and leave the other fields empty.
*/
var findResponse = struct {
	Results []core.Charm `json:"results"`
}{
	Results: []core.Charm{
		{
			ID:   "demo-id",
			Type: "charm",
			Name: "demo-charm",
			// Result, DefaultArtifact and Artifacts can be zero‑valued – they are
			// never examined in the tests that only verify the name.
			Result:          core.CharmResult{},
			DefaultArtifact: core.CharmArtifact{},
			Artifacts:       []core.CharmArtifact{},
		},
	},
}

/*
`info` endpoint – returns a full core.Charm object that contains an

	Artifacts slice.  Only the fields required for the test are filled.
*/
var infoResponse = core.Charm{
	Name: "demo-charm",
	Artifacts: []core.CharmArtifact{
		{
			Channel: core.CharmChannel{
				Name: "stable",
				Base: core.CharmBase{
					Architecture: "amd64",
					Channel:      "stable",
					Name:         "ubuntu",
				},
				Risk: "stable",
			},
			Revision: core.CharmRevision{
				Revision:  1,
				Version:   "1.0.0",
				CreatedAt: time.Now(),
				Bases: []core.CharmBase{
					{Architecture: "amd64", Channel: "stable", Name: "ubuntu"},
				},
			},
		},
	},
}

/* -------------------------------------------------------------------------- *
 * Construction / interface compliance
 * -------------------------------------------------------------------------- */
func TestNewCharm(t *testing.T) {
	j := charm_mustNewJuju(t, &config.Config{})
	r := NewCharm(j)
	if r == nil {
		t.Fatal("NewCharm returned nil")
	}
	var _ core.CharmRepo = r
}

/* -------------------------------------------------------------------------- *
 * Structure validation
 * -------------------------------------------------------------------------- */
func TestCharm_Structure(t *testing.T) {
	j := charm_mustNewJuju(t, &config.Config{})
	r := NewCharm(j)

	c, ok := r.(*charm)
	if !ok {
		t.Fatal("expected *charm, got a different type")
	}
	if c.juju == nil {
		t.Error("expected juju field to be set, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Successful List – HTTP server returns a well‑formed find response
 * -------------------------------------------------------------------------- */
func TestCharm_List_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/find") {
			_ = json.NewEncoder(w).Encode(findResponse)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Juju: config.Juju{
			CharmhubAPIURL: ts.URL,
		},
	}
	j := charm_mustNewJuju(t, cfg)
	repo := NewCharm(j)

	ctx := context.Background()
	ch, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(ch) == 0 {
		t.Fatalf("expected at least one charm, got none")
	}
	if ch[0].Name != "demo-charm" {
		t.Fatalf("expected charm name %q, got %q", "demo-charm", ch[0].Name)
	}
}

/* -------------------------------------------------------------------------- *
 * Get – charm present
 * -------------------------------------------------------------------------- */
func TestCharm_Get_Found(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// both the find and info endpoints are hit – we answer with the same payload.
		if strings.Contains(r.URL.Path, "/find") || strings.Contains(r.URL.Path, "/info") {
			_ = json.NewEncoder(w).Encode(findResponse)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Juju: config.Juju{
			CharmhubAPIURL: ts.URL,
		},
	}
	j := charm_mustNewJuju(t, cfg)
	repo := NewCharm(j)

	ctx := context.Background()
	ch, err := repo.Get(ctx, "demo-charm")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if ch.Name != "demo-charm" {
		t.Fatalf("expected charm name %q, got %q", "demo-charm", ch.Name)
	}
}

/* -------------------------------------------------------------------------- *
 * Get – charm not found (empty result set)
 * -------------------------------------------------------------------------- */
func TestCharm_Get_NotFound(t *testing.T) {
	// Return an empty result list – the code should translate this into a
	// connect.CodeNotFound error.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(struct {
			Results []core.Charm `json:"results"`
		}{})
	}))
	defer ts.Close()

	cfg := &config.Config{
		Juju: config.Juju{
			CharmhubAPIURL: ts.URL,
		},
	}
	j := charm_mustNewJuju(t, cfg)
	repo := NewCharm(j)

	ctx := context.Background()
	_, err := repo.Get(ctx, "missing")
	if err == nil {
		t.Fatalf("expected error for missing charm, got nil")
	}
	if connect.CodeOf(err) != connect.CodeNotFound {
		t.Fatalf("expected CodeNotFound, got %v", err)
	}
}

/* -------------------------------------------------------------------------- *
 * ListArtifacts – calls the “info” endpoint
 * -------------------------------------------------------------------------- */
func TestCharm_ListArtifacts_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/info/") {
			_ = json.NewEncoder(w).Encode(infoResponse)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Juju: config.Juju{
			CharmhubAPIURL: ts.URL,
		},
	}
	j := charm_mustNewJuju(t, cfg)
	repo := NewCharm(j)

	ctx := context.Background()
	art, err := repo.ListArtifacts(ctx, "demo-charm")
	if err != nil {
		t.Fatalf("ListArtifacts returned error: %v", err)
	}
	if len(art) == 0 {
		t.Fatalf("expected at least one artifact, got none")
	}
	if art[0].Channel.Name != "stable" {
		t.Fatalf("expected artifact channel name %q, got %q", "stable", art[0].Channel.Name)
	}
}

/* -------------------------------------------------------------------------- *
 * Error handling – malformed base URL
 * -------------------------------------------------------------------------- */
func TestCharm_Error_MalformedURL(t *testing.T) {
	cfg := &config.Config{
		Juju: config.Juju{
			CharmhubAPIURL: "::::",
		},
	}
	j := charm_mustNewJuju(t, cfg)
	repo := NewCharm(j)

	ctx := context.Background()
	_, err := repo.List(ctx)
	if err == nil {
		t.Fatalf("expected error when base URL is malformed, got nil")
	}
}

/* -------------------------------------------------------------------------- *
 * Edge case – nil Juju pointer should panic on any method call
 * -------------------------------------------------------------------------- */
func TestCharm_EdgeCases_NilJuju(t *testing.T) {
	r := &charm{juju: nil}
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when List is called with nil Juju")
			}
		}()
		_, _ = r.List(ctx)
	})

	t.Run("Get", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when Get is called with nil Juju")
			}
		}()
		_, _ = r.Get(ctx, "any")
	})

	t.Run("ListArtifacts", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when ListArtifacts is called with nil Juju")
			}
		}()
		_, _ = r.ListArtifacts(ctx, "any")
	})
}

/* -------------------------------------------------------------------------- *
 * Benchmark – creation only (method calls would need a real server)
 * -------------------------------------------------------------------------- */
func BenchmarkCharm_Creation(b *testing.B) {
	cfg := &config.Config{}
	j := charm_mustNewJuju(b, cfg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if r := NewCharm(j); r == nil {
			b.Fatal("creation returned nil")
		}
	}
}

/* -------------------------------------------------------------------------- *
 * Keep the utils import alive – the production code uses it.
 * -------------------------------------------------------------------------- */
var _ = wrap.HTTPGet // silence unused‑import warning
