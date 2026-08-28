package chisel

import (
	"context"
	"net"
	"strconv"
	"time"
)

const (
	// healthCheckInterval is how often every registered endpoint is probed.
	healthCheckInterval = 15 * time.Second

	// healthDialTimeout bounds a single TCP probe.
	healthDialTimeout = 2 * time.Second

	// healthFailThreshold is how many consecutive failures deregister a cluster.
	healthFailThreshold = 3
)

// HealthCheckListener wraps the health check loop as a transport.Listener, so
// it shares the errgroup lifecycle — and the panic handling and coordinated
// shutdown that come with it — with the HTTP and tunnel servers.
type HealthCheckListener struct {
	service *Service
}

func NewHealthCheckListener(service *Service) *HealthCheckListener {
	return &HealthCheckListener{service: service}
}

func (h *HealthCheckListener) Start(ctx context.Context) error {
	h.service.runHealthCheck(ctx)
	return nil
}

// Stop is a no-op: the loop exits when its context is canceled.
func (h *HealthCheckListener) Stop(_ context.Context) error {
	return nil
}

// clusterSnapshot copies the cluster-to-host mapping, so health checks can
// iterate without holding the lock.
func (s *Service) clusterSnapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make(map[string]string, len(s.links))
	for name, entry := range s.links {
		snapshot[name] = entry.Host
	}
	return snapshot
}

// runHealthCheck TCP-dials every registered endpoint, deregistering clusters
// that fail healthFailThreshold consecutive probes. It blocks until ctx is
// canceled.
func (s *Service) runHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	dialer := &net.Dialer{Timeout: healthDialTimeout}
	failCounts := make(map[string]int)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkClusters(ctx, dialer, failCounts)
		}
	}
}

// checkClusters runs one round, mutating failCounts in place.
func (s *Service) checkClusters(ctx context.Context, dialer *net.Dialer, failCounts map[string]int) {
	snapshot := s.clusterSnapshot()

	// Drop counts for clusters that are no longer registered.
	for name := range failCounts {
		if _, ok := snapshot[name]; !ok {
			delete(failCounts, name)
		}
	}

	for cluster, host := range snapshot {
		addr := net.JoinHostPort(host, strconv.Itoa(tunnelPort))
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			if closeErr := conn.Close(); closeErr != nil {
				s.log.Debug("failed to close health check connection", "cluster", cluster, "error", closeErr)
			}
			if failCounts[cluster] > 0 {
				s.log.Debug("cluster recovered", "cluster", cluster)
			}
			delete(failCounts, cluster)
			continue
		}

		// Cancellation is not a probe failure.
		if ctx.Err() != nil {
			return
		}

		failCounts[cluster]++
		s.log.Debug("probe failed",
			"cluster", cluster,
			"address", addr,
			"consecutive_failures", failCounts[cluster],
			"error", err,
		)

		if failCounts[cluster] >= healthFailThreshold {
			// A concurrent re-registration would have assigned a new host,
			// and deregistering that would be wrong.
			s.mu.RLock()
			current, exists := s.links[cluster]
			s.mu.RUnlock()
			if exists && current.Host == host {
				s.log.Info("deregistering disconnected cluster",
					"cluster", cluster,
					"consecutive_failures", failCounts[cluster],
				)
				s.DeregisterCluster(cluster)
			}
			delete(failCounts, cluster)
		}
	}
}
