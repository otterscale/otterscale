package handler

import (
	"fmt"
	"slices"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

// What the shared scope may hand to any Workspace member. Changing these
// lists, or the checks in sharedQueryProblems, needs a security review: they
// are what keeps a shared query from naming a namespace or a pod
// (OTTERSCALE-288, 2026-10-13-shared-metrics-design.md §5).
//
// Checked against HAMi v2.10.0 (phison-ai/release SXMR0105/0003) and the DCGM
// exporter on the HCI lab cluster, 2026-10-13.
var (
	// sharedSafeGrouping: per card, per host, per GPU index.
	sharedSafeGrouping = map[string]bool{
		"UUID": true, "modelName": true, "hostname": true, "device": true,
		"GPU_I_ID": true, "gpu": true, "device_uuid": true,
	}
	// sharedSafeMatchers: the labels the hub fills in from uuid/hostname.
	sharedSafeMatchers = map[string]bool{"UUID": true, "hostname": true}
	// sharedSafeMetrics: the metrics the catalog reads.
	sharedSafeMetrics = map[string]bool{
		"DCGM_FI_DEV_GPU_UTIL": true, "DCGM_FI_DEV_FB_USED": true, "DCGM_FI_DEV_FB_FREE": true,
		"DCGM_FI_DEV_FB_RESERVED": true, "DCGM_FI_DEV_XID_ERRORS": true,
		"hami_vgpu_memory_limit_bytes": true, "hami_vgpu_memory_used_bytes": true,
		"hami_gpu_memory_allocated_bytes": true, "hami_gpu_core_allocated_ratio": true,
		"hami_gpu_shared_count": true,
		"kubevirt_info":         true, "kubevirt_virt_controller_ready_status": true,
		"ceph_cluster_total_bytes": true,
	}
	// sharedSafeAggregations: aggregations whose output carries only their
	// grouping labels. topk, bottomk and limitk return their input series
	// whole; count_values turns values into a label.
	sharedSafeAggregations = map[parser.ItemType]bool{
		parser.SUM: true, parser.MAX: true, parser.MIN: true, parser.AVG: true, parser.COUNT: true,
	}
	// sharedForbiddenFunctions write labels from other labels.
	sharedForbiddenFunctions = map[string]bool{"label_replace": true, "label_join": true, "info": true}
)

// sharedQueryProblems lists every way promql could answer, or fail, with more
// than per-card, per-host or total numbers.
func sharedQueryProblems(promql string) []string {
	expr, err := parser.NewParser(parser.Options{}).ParseExpr(promql)
	if err != nil {
		return []string{err.Error()}
	}
	var problems []string
	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.AggregateExpr:
			if !sharedSafeAggregations[n.Op] {
				problems = append(problems, fmt.Sprintf("aggregation %s keeps or invents labels", n.Op))
			}
			if n.Without {
				problems = append(problems, "without keeps every label it does not name")
			}
			for _, label := range n.Grouping {
				if !sharedSafeGrouping[label] {
					problems = append(problems, fmt.Sprintf("groups by %q", label))
				}
			}
		case *parser.BinaryExpr:
			// "ignoring () group_left" names no label and is caught by its
			// card alone. Set operators are many-to-many by definition.
			vm := n.VectorMatching
			grouped := vm != nil && !n.Op.IsSetOperator() && vm.Card != parser.CardOneToOne
			if vm != nil && (vm.On || len(vm.MatchingLabels) > 0 || len(vm.Include) > 0 || grouped) {
				problems = append(problems,
					"explicit vector matching (on, ignoring, group_left/right) can name series in Prometheus' error")
			}
		case *parser.Call:
			if sharedForbiddenFunctions[n.Func.Name] {
				problems = append(problems, "calls "+n.Func.Name)
			}
		case *parser.VectorSelector:
			if !sharedSafeMetrics[n.Name] {
				problems = append(problems, fmt.Sprintf("reads %s", n))
			}
			problems = append(problems, sharedMatcherProblems(n)...)
			inAggregation := slices.ContainsFunc(path, func(p parser.Node) bool {
				_, ok := p.(*parser.AggregateExpr)
				return ok
			})
			if !inAggregation {
				problems = append(problems, fmt.Sprintf("%s is not inside an aggregation", n))
			}
		}
		return nil
	})
	return problems
}

// sharedMatcherProblems, part of sharedQueryProblems: a selector may match on
// its own metric name and on the labels the hub fills in, nothing else.
func sharedMatcherProblems(n *parser.VectorSelector) []string {
	var problems []string
	for _, m := range n.LabelMatchers {
		ownName := m.Name == "__name__" && m.Type == labels.MatchEqual && m.Value == n.Name
		if !ownName && !sharedSafeMatchers[m.Name] {
			problems = append(problems, fmt.Sprintf("%s filters on %s", n, m.Name))
		}
	}
	return problems
}

// TestSharedCatalog_Safety: every query, filled in, answers per card, per
// host or in total — never with a namespace, a pod or another workload label.
func TestSharedCatalog_Safety(t *testing.T) {
	for name, q := range sharedCatalog {
		t.Run(name, func(t *testing.T) {
			promql, err := buildSharedQuery(q, []string{"x"})
			if err != nil {
				t.Fatalf("build: %v", err)
			}
			for _, problem := range sharedQueryProblems(promql) {
				t.Errorf("%s: %s", promql, problem)
			}
		})
	}
}

// TestSharedQueryProblems_CatchesLeaks: the check itself refuses what would
// hand out workload identity.
func TestSharedQueryProblems_CatchesLeaks(t *testing.T) {
	for _, promql := range []string{
		`hami_vgpu_memory_used_bytes`,
		`sum without (device_uuid) (hami_vgpu_memory_used_bytes)`,
		`sum by (exported_namespace) (hami_vgpu_memory_used_bytes)`,
		`topk by (device_uuid) (1, hami_vgpu_memory_used_bytes)`,
		`count_values by (device_uuid) ("v", hami_vgpu_memory_used_bytes)`,
		`sum by (device_uuid) (label_replace(hami_vgpu_memory_used_bytes, "device_uuid", "$1", "exported_pod", "(.*)"))`,
		`sum by (UUID) (kube_pod_info)`,
		`sum by (UUID) ({__name__="DCGM_FI_DEV_FB_USED"})`,
		`sum(DCGM_FI_DEV_FB_USED) + DCGM_FI_DEV_FB_FREE`,
		`sum by (UUID) (DCGM_FI_DEV_FB_USED{exported_namespace="team-a"})`,
		`sum by (device_uuid) (hami_vgpu_memory_used_bytes * on (device_uuid) hami_vgpu_memory_limit_bytes)`,
		`sum by (device_uuid) (hami_vgpu_memory_used_bytes * ignoring () group_left hami_vgpu_memory_limit_bytes)`,
		`sum by (device_uuid) (hami_vgpu_memory_used_bytes > ignoring () group_right hami_vgpu_memory_limit_bytes)`,
	} {
		if _, err := parser.NewParser(parser.Options{}).ParseExpr(promql); err != nil {
			t.Errorf("%s does not parse: %v", promql, err)
			continue
		}
		if problems := sharedQueryProblems(promql); len(problems) == 0 {
			t.Errorf("%s passed the safety check", promql)
		}
	}
}

// TestSharedCatalog_Golden: the exact PromQL each name sends. The first
// fifteen are the queries the dashboard sent at cluster scope before; a change
// here changes what admins see too.
func TestSharedCatalog_Golden(t *testing.T) {
	const (
		u = `{UUID=~"GPU-3fd1\\.a|GPU-b"}`
		h = `{hostname=~"hci|gpu-node\\.lab"}`
	)
	values := map[sharedParam][]string{
		sharedUUID:     {"GPU-3fd1.a", "GPU-b"},
		sharedHostname: {"hci", "gpu-node.lab"},
	}
	want := map[string]string{
		"gpu-util":    `max by (UUID) (DCGM_FI_DEV_GPU_UTIL` + u + `)`,
		"gpu-fb-used": `sum by (UUID) (DCGM_FI_DEV_FB_USED` + u + `)`,
		"gpu-fb-total": `sum by (UUID) (DCGM_FI_DEV_FB_USED` + u + ` + DCGM_FI_DEV_FB_FREE` + u +
			` + DCGM_FI_DEV_FB_RESERVED` + u + `) or sum by (UUID) (DCGM_FI_DEV_FB_USED` + u +
			` + DCGM_FI_DEV_FB_FREE` + u + `)`,
		"gpu-inventory":              `sum by (UUID, modelName, hostname, device, GPU_I_ID) (DCGM_FI_DEV_FB_FREE + DCGM_FI_DEV_FB_RESERVED + DCGM_FI_DEV_FB_USED) * 1048576`,
		"vgpu-memory-limit":          `sum by (device_uuid) (hami_vgpu_memory_limit_bytes)`,
		"vgpu-memory-used":           `sum by (device_uuid) (hami_vgpu_memory_used_bytes)`,
		"host-gpu-fb-used":           `sum by (hostname) (DCGM_FI_DEV_FB_USED` + h + `)`,
		"host-gpu-fb-total":          `sum by (hostname) (DCGM_FI_DEV_FB_USED` + h + ` + DCGM_FI_DEV_FB_FREE` + h + `)`,
		"host-gpu-count":             `count by (hostname) (DCGM_FI_DEV_FB_USED` + h + `)`,
		"gpu-xid-total":              `sum(increase(DCGM_FI_DEV_XID_ERRORS` + h + `[24h]))`,
		"gpu-xid-by-gpu":             `sum by (hostname, gpu) (increase(DCGM_FI_DEV_XID_ERRORS` + h + `[24h])) > 0`,
		"kubevirt-handlers":          `count(kubevirt_info)`,
		"kubevirt-controllers":       `count(kubevirt_virt_controller_ready_status)`,
		"kubevirt-controllers-ready": `sum(kubevirt_virt_controller_ready_status)`,
		"ceph-total-bytes":           `sum(ceph_cluster_total_bytes) or vector(0)`,
		"gpu-allocated-memory":       `max by (device_uuid) (hami_gpu_memory_allocated_bytes)`,
		"gpu-allocated-cores":        `max by (device_uuid) (hami_gpu_core_allocated_ratio)`,
		"gpu-shared-count":           `max by (device_uuid) (hami_gpu_shared_count)`,
	}
	if len(sharedCatalog) != len(want) {
		t.Errorf("catalog has %d queries, want %d", len(sharedCatalog), len(want))
	}
	for name, q := range sharedCatalog {
		got, err := buildSharedQuery(q, values[q.param])
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if got != want[name] {
			t.Errorf("%s:\n got %s\nwant %s", name, got, want[name])
		}
		wantEndpoint := "/api/v1/query"
		if name == "host-gpu-fb-used" {
			wantEndpoint = "/api/v1/query_range"
		}
		if q.endpoint() != wantEndpoint {
			t.Errorf("%s: endpoint = %s, want %s", name, q.endpoint(), wantEndpoint)
		}
	}
}

// TestSharedMatcher_MatchesValuesLiterally: a value is matched as written,
// whatever regexp or PromQL syntax it contains.
func TestSharedMatcher_MatchesValuesLiterally(t *testing.T) {
	values := []string{`a"b|c.d`, "gpu-node.lab"}
	ms, err := parser.NewParser(parser.Options{}).ParseMetricSelector("{" + sharedMatcher("hostname", values) + "}")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(ms) != 1 {
		t.Fatalf("matchers = %v, want one", ms)
	}
	for _, v := range values {
		if !ms[0].Matches(v) {
			t.Errorf("does not match %q", v)
		}
	}
	for _, v := range []string{`a"b`, "c.d", "gpu-nodexlab", "gpu-node.lab2", "xgpu-node.lab"} {
		if ms[0].Matches(v) {
			t.Errorf("also matches %q", v)
		}
	}
}
