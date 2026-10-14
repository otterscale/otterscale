package handler

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	// sharedMaxValues caps each repeatable parameter: one request names at
	// most this many cards or hosts.
	sharedMaxValues = 256

	timeParam  = "time"
	startParam = "start"
	endParam   = "end"
	stepParam  = "step"
)

// sharedUUIDPattern is a GPU UUID as DCGM and HAMi report it; the dashboard
// applies the same rule before asking (fetch-utilization.ts, SAFE_UUID).
var sharedUUIDPattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,128}$`)

var errBadSharedRequest = errors.New("bad request")

// ServeShared answers /proxy/{cluster}/prometheus-shared/{name} with one of
// the hub's own named queries (sharedCatalog, proxy_shared_catalog.go):
// cluster-wide numbers that name no workload — per GPU card, per host, or a
// total — for Workspace members as well as cluster viewers. Callers send a
// name and parameters, never PromQL.
//
// namespace, when given, only proves that the caller belongs to a Workspace:
// it is checked like the namespace scope's (get pods there) and goes no
// further, since the answer is not per namespace. A SelfSubjectAccessReview
// asks about one namespace at a time, never "any Workspace", hence the
// parameter. Without it, only cluster viewers get through.
func (h *ProxyHandler) ServeShared(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writePrometheusError(w, http.StatusMethodNotAllowed, "bad_data", "only GET")
		return
	}
	cluster := r.PathValue("cluster")
	q, ok := sharedCatalog[r.PathValue("name")]
	if !ok {
		writePrometheusError(w, http.StatusNotFound, "not_found", "no shared query by that name")
		return
	}

	params, namespace, err := parseSharedRequest(r, q)
	if err != nil {
		if errors.Is(err, errSharedTemplate) {
			writePrometheusError(w, http.StatusInternalServerError, "internal", err.Error())
		} else {
			writePrometheusError(w, http.StatusBadRequest, "bad_data", err.Error())
		}
		return
	}

	target, ok := h.resolve(w, r, cluster)
	if !ok {
		return
	}
	attrs := []core.AccessAttributes{core.ClusterMetricsAccess}
	if namespace != "" {
		attrs = append(attrs, core.NamespaceMetricsAccess(namespace))
	}
	if !h.authorize(w, r, cluster, attrs...) {
		return
	}

	r.Body = http.NoBody
	r.ContentLength = 0
	h.forward(w, r, target, q.endpoint(), params.Encode())
}

// parseSharedRequest checks every parameter against what q declares and
// returns what Prometheus gets: the filled-in query, plus time, or start, end
// and step.
func parseSharedRequest(r *http.Request, q sharedQuery) (params url.Values, namespace string, err error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, "", fmt.Errorf("%w: query string: %w", errBadSharedRequest, err)
	}

	timing := []string{timeParam}
	if q.isRange {
		timing = []string{startParam, endParam, stepParam}
	}
	allowed := map[string]bool{namespaceParam: true}
	if q.param != sharedNoParam {
		allowed[string(q.param)] = true
	}
	for _, p := range timing {
		allowed[p] = true
	}
	for key := range query {
		if !allowed[key] {
			return nil, "", fmt.Errorf("%w: unknown parameter %q", errBadSharedRequest, key)
		}
	}

	if values, ok := query[namespaceParam]; ok {
		if len(values) != 1 {
			return nil, "", fmt.Errorf("%w: at most one namespace parameter", errBadSharedRequest)
		}
		namespace = values[0]
		if problems := validation.IsDNS1123Label(namespace); len(problems) > 0 {
			return nil, "", fmt.Errorf("%w: namespace %q: %s", errBadSharedRequest, namespace, strings.Join(problems, "; "))
		}
	}

	values, err := sharedValues(q.param, query[string(q.param)])
	if err != nil {
		return nil, "", err
	}
	promql, err := buildSharedQuery(q, values)
	if err != nil {
		return nil, "", err
	}

	params = url.Values{queryParam: {promql}}
	for _, p := range timing {
		given := query[p]
		switch {
		case len(given) == 0 && q.isRange:
			return nil, "", fmt.Errorf("%w: %s is required", errBadSharedRequest, p)
		case len(given) == 0:
			continue
		case len(given) > 1:
			return nil, "", fmt.Errorf("%w: %s given more than once", errBadSharedRequest, p)
		case !validSharedSeconds(p, given[0]):
			return nil, "", fmt.Errorf("%w: %s %q is not a number of seconds", errBadSharedRequest, p, given[0])
		}
		params.Set(p, given[0])
	}
	return params, namespace, nil
}

// sharedValues checks the values of a query's repeatable parameter.
func sharedValues(param sharedParam, values []string) ([]string, error) {
	if param == sharedNoParam {
		return nil, nil
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("%w: %s is required", errBadSharedRequest, param)
	}
	if len(values) > sharedMaxValues {
		return nil, fmt.Errorf("%w: at most %d %s values", errBadSharedRequest, sharedMaxValues, param)
	}
	for _, v := range values {
		if !validSharedValue(param, v) {
			return nil, fmt.Errorf("%w: %s %q", errBadSharedRequest, param, v)
		}
	}
	return values, nil
}

func validSharedValue(param sharedParam, v string) bool {
	switch param {
	case sharedUUID:
		return sharedUUIDPattern.MatchString(v)
	case sharedHostname:
		return len(validation.IsDNS1123Subdomain(v)) == 0
	case sharedNoParam:
		return false
	}
	return false
}

// validSharedSeconds: time, start and end are Unix seconds; step is seconds
// and more than zero.
func validSharedSeconds(param, v string) bool {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return false
	}
	return param != stepParam || f > 0
}
