package integration

import (
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
	"github.com/otterscale/otterscale/internal/providers/chisel"
	tunneltransport "github.com/otterscale/otterscale/internal/transport/tunnel"
)

// integrationSecret backs the enrolment tokens these tests present.
const integrationSecret = "integration-root-secret"

// newTestLink builds a LinkUseCase with enrolment configured, and a
// helper that mints the token for a cluster.
func newTestLink(t *testing.T, tunnel core.TunnelProvider) (link *core.LinkUseCase, tokenFor func(cluster string) string) {
	t.Helper()

	enrolment, err := core.NewEnrolment(integrationSecret)
	if err != nil {
		t.Fatalf("NewEnrolment: %v", err)
	}
	return core.NewLinkUseCase(tunnel, "test", enrolment), enrolment.Token
}

func TestLinkRegisterClusterUsesSingleSharedTunnelPort(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	csrA := generateCSR(t, "agent-a")
	csrB := generateCSR(t, "agent-b")

	regA, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-a",
		AgentID:        "agent-a",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-a"),
		CSRPEM:         csrA,
	})
	if err != nil {
		t.Fatalf("register cluster-a: %v", err)
	}
	regB, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-b",
		AgentID:        "agent-b",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-b"),
		CSRPEM:         csrB,
	})
	if err != nil {
		t.Fatalf("register cluster-b: %v", err)
	}

	if len(regA.Certificate) == 0 || len(regB.Certificate) == 0 {
		t.Fatal("expected non-empty certificates")
	}
	if len(regA.CACertificate) == 0 || len(regB.CACertificate) == 0 {
		t.Fatal("expected non-empty CA certificates")
	}

	if regA.Endpoint == "" || regB.Endpoint == "" {
		t.Fatalf("expected non-empty tunnel endpoints, got endpointA=%q endpointB=%q", regA.Endpoint, regB.Endpoint)
	}
	if regA.Endpoint == regB.Endpoint {
		t.Fatalf("expected distinct endpoints for different clusters, got %q", regA.Endpoint)
	}

	addrA, err := tunnel.ResolveAddress(t.Context(), "cluster-a")
	if err != nil {
		t.Fatalf("resolve cluster-a: %v", err)
	}
	addrB, err := tunnel.ResolveAddress(t.Context(), "cluster-b")
	if err != nil {
		t.Fatalf("resolve cluster-b: %v", err)
	}

	if !strings.HasSuffix(addrA, ":16598") || !strings.HasSuffix(addrB, ":16598") {
		t.Fatalf("expected resolved addresses to use shared port 16598, got addrA=%q addrB=%q", addrA, addrB)
	}
}

func TestLinkRegisterClusterLatestAgentWinsForSameCluster(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	csr1 := generateCSR(t, "agent-r-1")
	csr2 := generateCSR(t, "agent-r-2")

	_, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-r",
		AgentID:        "agent-r-1",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-r"),
		CSRPEM:         csr1,
	})
	if err != nil {
		t.Fatalf("register agent-r-1: %v", err)
	}
	reg2, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-r",
		AgentID:        "agent-r-2",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-r"),
		CSRPEM:         csr2,
	})
	if err != nil {
		t.Fatalf("register agent-r-2: %v", err)
	}

	// After re-registration the route must resolve to the latest
	// agent's endpoint regardless of whether the host was reused.
	addr, err := tunnel.ResolveAddress(t.Context(), "cluster-r")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if addr != "http://"+reg2.Endpoint {
		t.Fatalf("expected resolve to use latest agent endpoint %q, got %q", reg2.Endpoint, addr)
	}

	// Only one cluster should be registered.
	links := tunnel.ListLinks()
	if len(links) != 1 || slices.Collect(maps.Keys(links))[0] != "cluster-r" {
		t.Fatalf("expected exactly one cluster 'cluster-r', got %v", links)
	}
}

func TestLinkRegisterClusterReregisterAndReplaceAcrossAgents(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	csrA := generateCSR(t, "agent-a")
	csrB := generateCSR(t, "agent-b")

	regA1, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-z",
		AgentID:        "agent-a",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-z"),
		CSRPEM:         csrA,
	})
	if err != nil {
		t.Fatalf("register agent-a #1: %v", err)
	}

	regB, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-z",
		AgentID:        "agent-b",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-z"),
		CSRPEM:         csrB,
	})
	if err != nil {
		t.Fatalf("register agent-b: %v", err)
	}

	// After re-registration for the same cluster, the route must
	// resolve to the latest agent's endpoint.
	addrB, err := tunnel.ResolveAddress(t.Context(), "cluster-z")
	if err != nil {
		t.Fatalf("resolve after agent-b register: %v", err)
	}
	if addrB != "http://"+regB.Endpoint {
		t.Fatalf("expected resolve to point to agent-b endpoint %q, got %q", regB.Endpoint, addrB)
	}

	regA2, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-z",
		AgentID:        "agent-a",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-z"),
		CSRPEM:         csrA,
	})
	if err != nil {
		t.Fatalf("register agent-a #2: %v", err)
	}

	// Every registration issues a fresh password, so re-registering
	// invalidates whatever the previous session was using.
	if regA1.TunnelPassword == "" || regA2.TunnelPassword == "" {
		t.Fatal("expected a tunnel password on every registration")
	}
	if regA1.TunnelPassword == regA2.TunnelPassword {
		t.Fatal("expected password rotation for same agent re-register")
	}

	// The tunnel user is the cluster, not the agent: agent-a and
	// agent-b claimed the same cluster under different agent ids and
	// must have been issued the same user name.
	if regA1.TunnelUser != "cluster-z" || regB.TunnelUser != "cluster-z" {
		t.Fatalf("tunnel users = %q / %q, want both %q", regA1.TunnelUser, regB.TunnelUser, "cluster-z")
	}

	for i := range 3 {
		addr, err := tunnel.ResolveAddress(t.Context(), "cluster-z")
		if err != nil {
			t.Fatalf("resolve #%d: %v", i+1, err)
		}
		if addr != "http://"+regA2.Endpoint {
			t.Fatalf("expected only re-registered route to be selected, got %q", addr)
		}
	}
}

// newTestTunnel creates a chisel.Service with a fresh test CA
// injected at construction time.
func newTestTunnel(t *testing.T) *chisel.Service {
	t.Helper()
	ca, err := pki.NewCA()
	if err != nil {
		t.Fatalf("create CA: %v", err)
	}
	return chisel.NewService(ca)
}

func initTunnelServer(t *testing.T, tunnel *chisel.Service) {
	t.Helper()

	srv, err := tunneltransport.NewServer(
		tunneltransport.WithServer(tunnel.ServerRef()),
	)
	if err != nil {
		t.Fatalf("init tunnel server: %v", err)
	}
	t.Cleanup(func() {
		_ = srv.Stop(t.Context())
	})
}

// generateCSR creates a fresh ECDSA key pair and PEM-encoded CSR for
// the given common name.
func generateCSR(t *testing.T, cn string) []byte {
	t.Helper()
	key, _, err := pki.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	csr, err := pki.GenerateCSR(key, cn)
	if err != nil {
		t.Fatalf("generate CSR: %v", err)
	}
	return csr
}

// TestLinkRegisterClusterRejectedTokenKeepsExistingAgent checks the
// property end to end, against the real tunnel provider: a rejected
// registration must leave the cluster pointing at the agent that is
// already serving it. Re-registration deletes the previous chisel user
// and releases its address, so a check that ran too late would let an
// unauthorized caller knock a healthy cluster offline.
func TestLinkRegisterClusterRejectedTokenKeepsExistingAgent(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	if _, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-x",
		AgentID:        "agent-x",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("cluster-x"),
		CSRPEM:         generateCSR(t, "agent-x"),
	}); err != nil {
		t.Fatalf("register the legitimate agent: %v", err)
	}

	before, err := tunnel.ResolveAddress(t.Context(), "cluster-x")
	if err != nil {
		t.Fatalf("resolve cluster-x: %v", err)
	}

	_, err = link.RegisterCluster(t.Context(), &core.RegistrationRequest{
		Cluster:        "cluster-x",
		AgentID:        "impostor",
		AgentVersion:   "test",
		EnrolmentToken: tokenFor("some-other-cluster"),
		CSRPEM:         generateCSR(t, "impostor"),
	})
	if err == nil {
		t.Fatal("registration with another cluster's token succeeded")
	}

	after, err := tunnel.ResolveAddress(t.Context(), "cluster-x")
	if err != nil {
		t.Fatalf("cluster-x is no longer registered after a rejected attempt: %v", err)
	}
	if after != before {
		t.Errorf("cluster-x moved from %q to %q after a rejected attempt", before, after)
	}
}

// TestLinkRegisterClusterSharedAgentHostnameDoesNotCollide is the
// regression test for cross-cluster credential clobbering.
//
// Agents identify themselves by hostname, and two clusters running the
// same deployment report the same one — an `otterscale-agent-0` in each
// is the ordinary case, not a contrived one. chisel keys its user index
// by name, so while the tunnel user was named after the agent, the
// second cluster's registration overwrote the first's credentials and
// stranded a working tunnel: the first agent could no longer
// authenticate, and the health checker eventually deregistered it.
func TestLinkRegisterClusterSharedAgentHostnameDoesNotCollide(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	// The hostname two identically-deployed agents both report.
	const sharedAgentID = "otterscale-agent-0"

	register := func(cluster string) core.Registration {
		t.Helper()
		reg, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
			Cluster:        cluster,
			AgentID:        sharedAgentID,
			AgentVersion:   "test",
			EnrolmentToken: tokenFor(cluster),
			CSRPEM:         generateCSR(t, sharedAgentID),
		})
		if err != nil {
			t.Fatalf("register %s: %v", cluster, err)
		}
		return reg
	}

	regA := register("cluster-a")
	regB := register("cluster-b")

	if regA.TunnelUser == "" || regB.TunnelUser == "" {
		t.Fatal("expected a tunnel user on every registration")
	}
	if regA.TunnelUser == regB.TunnelUser {
		t.Fatalf("both clusters were issued tunnel user %q, so one registration overwrites the other", regA.TunnelUser)
	}
	if regA.TunnelPassword == regB.TunnelPassword {
		t.Fatal("expected distinct tunnel passwords for distinct clusters")
	}

	// The first cluster must be untouched by the second's arrival.
	links := tunnel.ListLinks()
	if got := links["cluster-a"].User; got != regA.TunnelUser {
		t.Errorf("registry records tunnel user %q for cluster-a, but the agent was issued %q", got, regA.TunnelUser)
	}
	if links["cluster-a"].Host == links["cluster-b"].Host {
		t.Errorf("expected distinct loopback hosts, both got %q", links["cluster-a"].Host)
	}
	if _, err := tunnel.ResolveAddress(t.Context(), "cluster-a"); err != nil {
		t.Errorf("cluster-a stopped resolving once cluster-b registered: %v", err)
	}
}

// TestDeregisterClusterLeavesOtherClustersRegistered covers the second
// half of the same bug. The registry records the tunnel user it issued,
// and deregistration deletes that name from chisel — so while two
// clusters could record the same name, deregistering one revoked the
// other's credentials.
func TestDeregisterClusterLeavesOtherClustersRegistered(t *testing.T) {
	tunnel := newTestTunnel(t)
	initTunnelServer(t, tunnel)
	link, tokenFor := newTestLink(t, tunnel)

	const sharedAgentID = "otterscale-agent-0"

	for _, cluster := range []string{"cluster-a", "cluster-b"} {
		if _, err := link.RegisterCluster(t.Context(), &core.RegistrationRequest{
			Cluster:        cluster,
			AgentID:        sharedAgentID,
			AgentVersion:   "test",
			EnrolmentToken: tokenFor(cluster),
			CSRPEM:         generateCSR(t, sharedAgentID),
		}); err != nil {
			t.Fatalf("register %s: %v", cluster, err)
		}
	}

	before := tunnel.ListLinks()["cluster-a"]

	tunnel.DeregisterCluster("cluster-b")

	links := tunnel.ListLinks()
	if _, ok := links["cluster-b"]; ok {
		t.Error("cluster-b is still registered after deregistration")
	}
	if links["cluster-a"] != before {
		t.Errorf("cluster-a changed from %+v to %+v when cluster-b was deregistered", before, links["cluster-a"])
	}
	if _, err := tunnel.ResolveAddress(t.Context(), "cluster-a"); err != nil {
		t.Errorf("cluster-a stopped resolving when cluster-b was deregistered: %v", err)
	}
}
