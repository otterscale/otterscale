package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/otterscale/otterscale/internal/core"
)

// ssarServer plays kube-apiserver's SelfSubjectAccessReview endpoint and
// records what it was asked and as whom.
type ssarServer struct {
	mu      sync.Mutex
	allowed bool
	verdict *authorizationv1.SubjectAccessReviewStatus // non-nil: answer with this instead of allowed
	status  int                                        // non-zero: answer with this HTTP status instead
	got     []authorizationv1.ResourceAttributes
	users   []string
	groups  [][]string
}

func (s *ssarServer) handler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		// client-go sends core types as protobuf; the universal deserializer
		// reads that and JSON alike.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read review: %v", err)
		}
		var review authorizationv1.SelfSubjectAccessReview
		if _, _, err := scheme.Codecs.UniversalDeserializer().Decode(body, nil, &review); err != nil {
			t.Errorf("decode review: %v", err)
		}
		s.mu.Lock()
		if review.Spec.ResourceAttributes != nil {
			s.got = append(s.got, *review.Spec.ResourceAttributes)
		}
		s.users = append(s.users, r.Header.Get("Impersonate-User"))
		s.groups = append(s.groups, r.Header.Values("Impersonate-Group"))
		status, allowed, verdict := s.status, s.allowed, s.verdict
		s.mu.Unlock()

		if status != 0 {
			w.WriteHeader(status)
			return
		}
		review.Status = authorizationv1.SubjectAccessReviewStatus{Allowed: allowed}
		if verdict != nil {
			review.Status = *verdict
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(review)
	})
}

func newAccessReviewTest(t *testing.T, srv *ssarServer) core.AccessReviewer {
	t.Helper()
	httpSrv := httptest.NewServer(srv.handler(t))
	t.Cleanup(httpSrv.Close)
	return NewAccessReviewer(New(fakeTunnel{address: httpSrv.URL}))
}

func aliceContext(t *testing.T) context.Context {
	return core.WithUserInfo(t.Context(), core.UserInfo{
		Subject: "alice",
		Groups:  []string{"system:authenticated", "oidc:team-a"},
	})
}

// TestAccessReviewerAsksAsTheCaller pins the whole point: the question goes to
// the cluster as the signed-in user, with exactly the attributes given.
func TestAccessReviewerAsksAsTheCaller(t *testing.T) {
	srv := &ssarServer{allowed: true}
	reviewer := newAccessReviewTest(t, srv)

	allowed, err := reviewer.Allowed(aliceContext(t), "prod", core.ClusterMetricsAccess)
	if err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	if !allowed {
		t.Error("allowed = false, want true")
	}

	want := authorizationv1.ResourceAttributes{
		Namespace: "monitoring", Verb: "get", Group: "monitoring.coreos.com",
		Resource: "prometheuses", Subresource: "api", Name: "prometheus-stack-kube-prom-prometheus",
	}
	if len(srv.got) != 1 || srv.got[0] != want {
		t.Errorf("review attributes = %+v, want [%+v]", srv.got, want)
	}
	if len(srv.users) != 1 || srv.users[0] != "alice" {
		t.Errorf("Impersonate-User = %v, want [alice]", srv.users)
	}
	if want := []string{"system:authenticated", "oidc:team-a"}; len(srv.groups) != 1 || !slices.Equal(srv.groups[0], want) {
		t.Errorf("Impersonate-Group = %v, want [%v]", srv.groups, want)
	}
}

// TestAccessReviewerExplicitNo: an authorizer that denies outright, or that
// erred while evaluating and did not allow, has answered "no" — not failed.
func TestAccessReviewerExplicitNo(t *testing.T) {
	verdicts := map[string]authorizationv1.SubjectAccessReviewStatus{
		"denied":                     {Allowed: false, Denied: true, Reason: "denied by policy"},
		"evaluation error, no allow": {Allowed: false, EvaluationError: "role \"gone\" not found"},
	}
	for name, verdict := range verdicts {
		t.Run(name, func(t *testing.T) {
			reviewer := newAccessReviewTest(t, &ssarServer{verdict: &verdict})

			allowed, err := reviewer.Allowed(aliceContext(t), "prod", core.ClusterMetricsAccess)
			if err != nil {
				t.Fatalf("Allowed: %v, want a plain no", err)
			}
			if allowed {
				t.Error("allowed = true, want false")
			}
		})
	}
}

func TestAccessReviewerDenied(t *testing.T) {
	reviewer := newAccessReviewTest(t, &ssarServer{allowed: false})

	allowed, err := reviewer.Allowed(aliceContext(t), "prod", core.NamespaceMetricsAccess("ws-a"))
	if err != nil {
		t.Fatalf("Allowed: %v", err)
	}
	if allowed {
		t.Error("allowed = true, want false")
	}
}

// TestAccessReviewerFailureIsAnError: a review the cluster could not answer
// must not come back as a plain "no" (or worse, "yes") — callers turn it into 503.
func TestAccessReviewerFailureIsAnError(t *testing.T) {
	reviewer := newAccessReviewTest(t, &ssarServer{status: http.StatusInternalServerError})

	if _, err := reviewer.Allowed(aliceContext(t), "prod", core.ClusterMetricsAccess); err == nil {
		t.Fatal("Allowed returned no error for a failed review")
	}
}

func TestAccessReviewerRequiresACaller(t *testing.T) {
	reviewer := newAccessReviewTest(t, &ssarServer{allowed: true})

	_, err := reviewer.Allowed(t.Context(), "prod", core.ClusterMetricsAccess)
	if code, ok := core.DomainErrorCode(err); !ok || code != core.ErrorCodeUnauthenticated {
		t.Errorf("error = %v, want ErrorCodeUnauthenticated", err)
	}
}

// missingTunnel knows no cluster.
type missingTunnel struct{ fakeTunnel }

func (missingTunnel) ResolveAddress(_ context.Context, cluster string) (string, error) {
	return "", &core.ErrClusterNotFound{Cluster: cluster}
}

func TestAccessReviewerUnknownCluster(t *testing.T) {
	reviewer := NewAccessReviewer(New(missingTunnel{}))

	_, err := reviewer.Allowed(aliceContext(t), "gone", core.ClusterMetricsAccess)
	var notFound *core.ErrClusterNotFound
	if !errors.As(err, &notFound) {
		t.Errorf("error = %v, want *core.ErrClusterNotFound", err)
	}
}
