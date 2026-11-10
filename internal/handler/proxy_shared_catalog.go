package handler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/prometheus/promql/parser"
)

// sharedParam names the repeatable parameter a shared query takes.
type sharedParam string

const (
	sharedNoParam  sharedParam = ""
	sharedUUID     sharedParam = "uuid"     // a GPU UUID, DCGM's UUID label
	sharedHostname sharedParam = "hostname" // a node name, DCGM's hostname label
)

// errSharedTemplate is a hub bug: a catalog query that does not parse once
// filled in.
var errSharedTemplate = errors.New("shared query does not parse")

// sharedQuery is one named query of the shared scope (proxy_shared.go). The
// PromQL is the hub's own: callers pick a name and give values for the one
// parameter it declares, never PromQL. {U} and {H} stand for the uuid and
// hostname matchers.
//
// What a query may answer is pinned by TestSharedCatalog_Safety: every
// selector inside an aggregation by allowlisted labels, so the answer is per
// card, per host or a total — never a namespace, a pod or any other
// workload's identity.
type sharedQuery struct {
	promql  string
	param   sharedParam
	isRange bool
}

// sharedCatalog is every shared query by name. The first fifteen are the
// queries the dashboard used to send at cluster scope, unchanged in what they
// answer; the last three are hami-scheduler's per-card allocation.
//
// Exporters differ between clusters. When a metric is renamed, join the old
// and the new spelling with `or` in the same query, as gpu-fb-total does for
// DCGM exporters without FB_RESERVED, so one hub serves both.
var sharedCatalog = map[string]sharedQuery{
	"gpu-util":    {promql: `max by (UUID) (DCGM_FI_DEV_GPU_UTIL{U})`, param: sharedUUID},
	"gpu-fb-used": {promql: `sum by (UUID) (DCGM_FI_DEV_FB_USED{U})`, param: sharedUUID},
	"gpu-fb-total": {
		promql: `sum by (UUID) (DCGM_FI_DEV_FB_USED{U} + DCGM_FI_DEV_FB_FREE{U} + DCGM_FI_DEV_FB_RESERVED{U})` +
			` or sum by (UUID) (DCGM_FI_DEV_FB_USED{U} + DCGM_FI_DEV_FB_FREE{U})`,
		param: sharedUUID,
	},
	"gpu-inventory": {
		promql: `sum by (UUID, modelName, hostname, device, GPU_I_ID) (DCGM_FI_DEV_FB_FREE + DCGM_FI_DEV_FB_RESERVED + DCGM_FI_DEV_FB_USED) * 1048576`,
	},
	"vgpu-memory-limit": {promql: `sum by (device_uuid) (hami_vgpu_memory_limit_bytes)`},
	"vgpu-memory-used":  {promql: `sum by (device_uuid) (hami_vgpu_memory_used_bytes)`},
	"host-gpu-fb-used": {
		promql:  `sum by (hostname) (DCGM_FI_DEV_FB_USED{H})`,
		param:   sharedHostname,
		isRange: true,
	},
	"host-gpu-fb-total": {
		promql: `sum by (hostname) (DCGM_FI_DEV_FB_USED{H} + DCGM_FI_DEV_FB_FREE{H})`,
		param:  sharedHostname,
	},
	"host-gpu-count": {promql: `count by (hostname) (DCGM_FI_DEV_FB_USED{H})`, param: sharedHostname},
	"gpu-xid-total": {
		promql: `sum(increase(DCGM_FI_DEV_XID_ERRORS{H}[24h]))`,
		param:  sharedHostname,
	},
	"gpu-xid-by-gpu": {
		promql: `sum by (hostname, gpu) (increase(DCGM_FI_DEV_XID_ERRORS{H}[24h])) > 0`,
		param:  sharedHostname,
	},
	"kubevirt-handlers":          {promql: `count(kubevirt_info)`},
	"kubevirt-controllers":       {promql: `count(kubevirt_virt_controller_ready_status)`},
	"kubevirt-controllers-ready": {promql: `sum(kubevirt_virt_controller_ready_status)`},
	"ceph-total-bytes":           {promql: `sum(ceph_cluster_total_bytes) or vector(0)`},
	// hami-scheduler: one series per card, labeled by card and node only.
	// max drops the duplicates of a scheduler with several replicas.
	"gpu-allocated-memory": {promql: `max by (device_uuid) (hami_gpu_memory_allocated_bytes)`},
	"gpu-allocated-cores":  {promql: `max by (device_uuid) (hami_gpu_core_allocated_ratio)`},
	"gpu-shared-count":     {promql: `max by (device_uuid) (hami_gpu_shared_count)`},
}

// endpoint is the Prometheus API path the query goes to.
func (q sharedQuery) endpoint() string {
	if q.isRange {
		return "/api/v1/query_range"
	}
	return "/api/v1/query"
}

// sharedMatcher is label=~"v1|v2|…" with each value matched literally:
// regexp syntax is escaped first, then the whole is quoted as a PromQL string.
// Prometheus anchors the regexp, so nothing else matches.
func sharedMatcher(label string, values []string) string {
	literal := make([]string, len(values))
	for i, v := range values {
		literal[i] = regexp.QuoteMeta(v)
	}
	return label + "=~" + strconv.Quote(strings.Join(literal, "|"))
}

// buildSharedQuery fills q's matchers with values and checks the result parses.
func buildSharedQuery(q sharedQuery, values []string) (string, error) {
	promql := q.promql
	switch q.param {
	case sharedUUID:
		promql = strings.ReplaceAll(promql, "{U}", "{"+sharedMatcher("UUID", values)+"}")
	case sharedHostname:
		promql = strings.ReplaceAll(promql, "{H}", "{"+sharedMatcher("hostname", values)+"}")
	case sharedNoParam:
		// nothing to fill in
	}
	if _, err := parser.NewParser(parser.Options{}).ParseExpr(promql); err != nil {
		return "", fmt.Errorf("%w: %w", errSharedTemplate, err)
	}
	return promql, nil
}
