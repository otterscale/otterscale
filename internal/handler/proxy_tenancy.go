package handler

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	// namespaceParam is read from the URL only, as OpenShift's kube-rbac-proxy
	// does it (rewrites.byQueryParameter), and never forwarded.
	namespaceParam = "namespace"
	namespaceLabel = "namespace"
	queryParam     = "query"
	matchParam     = "match[]"

	formContentType = "application/x-www-form-urlencoded"
	// maxFormBytes is net/http's own ceiling for a parsed form body, which
	// Prometheus' form parsing has too. It applies to both the caller's body
	// and the body the hub forwards.
	maxFormBytes = 10 << 20
)

// tenancyKind is what namespace scope does with an endpoint.
type tenancyKind int

const (
	tenancyForbidden tenancyKind = iota
	tenancyQuery                 // PromQL in `query`
	tenancyMatch                 // series selectors in `match[]`
)

var (
	errBadTenancyRequest = errors.New("bad request")
	errBodyTooLarge      = errors.New("request body too large")
	// errNoNamespaceRule is a hub bug: an endpoint kind reached the rewrite
	// without a rule for it.
	errNoNamespaceRule = errors.New("no namespace rule for this endpoint")
)

// ServeTenancy forwards /proxy/{cluster}/prometheus-tenancy/{path...}?namespace=<ns>
// for callers who may get pods in <ns>, or who have cluster metrics access,
// rewritten to namespace="<ns>" by prom-label-proxy's enforcer — the package
// OpenShift's prom-label-proxy sidecar runs, with its --error-on-replace: a
// selector that cannot match <ns> is refused, one that can is narrowed to it.
// The rewrite is the enforcer's, unchanged; the hub adds no rules of its own.
//
// What is forwarded is always the hub's own re-encoding of what it checked,
// never the request's raw query string or body, so Prometheus cannot read a
// parameter the hub did not see.
func (h *ProxyHandler) ServeTenancy(w http.ResponseWriter, r *http.Request) {
	cluster := r.PathValue("cluster")

	promPath, allowed := core.AllowedPrometheusPath("/" + r.PathValue("path"))
	kind := tenancyEndpoint(promPath)
	if !allowed || kind == tenancyForbidden {
		writePrometheusError(w, http.StatusForbidden, "forbidden", "endpoint not available at namespace scope")
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		writePrometheusError(w, http.StatusMethodNotAllowed, "bad_data", "only GET and POST")
		return
	}

	query, form, namespace, err := parseTenancyRequest(w, r)
	if err == nil {
		err = enforceNamespace(kind, namespace, query, form)
	}
	var body string
	if err == nil && r.Method == http.MethodPost {
		body, err = forwardedForm(form)
	}
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, errBodyTooLarge):
			status = http.StatusRequestEntityTooLarge
		case errors.Is(err, injectproxy.ErrEnforceLabel), errors.Is(err, errNoNamespaceRule):
			status = http.StatusInternalServerError
		}
		writePrometheusError(w, status, "bad_data", err.Error())
		return
	}

	target, ok := h.resolve(w, r, cluster)
	if !ok {
		return
	}
	if !h.authorize(w, r, cluster, core.ClusterMetricsAccess, core.NamespaceMetricsAccess(namespace)) {
		return
	}

	contentType := ""
	if r.Method == http.MethodPost {
		r.Body = io.NopCloser(strings.NewReader(body))
		r.ContentLength = int64(len(body))
		contentType = formContentType
	} else {
		r.Body = http.NoBody
		r.ContentLength = 0
	}
	// The body above is the hub's own plain re-encoding; the caller's
	// Content-Encoding described a body that is not being forwarded.
	r.Header.Del("Content-Encoding")
	h.forward(w, r, target, promPath, query.Encode(), contentType)
}

func tenancyEndpoint(promPath string) tenancyKind {
	switch promPath {
	case "/api/v1/query", "/api/v1/query_range":
		return tenancyQuery
	case "/api/v1/series", "/api/v1/labels":
		return tenancyMatch
	}
	if name, ok := strings.CutPrefix(promPath, "/api/v1/label/"); ok {
		if name, ok = strings.CutSuffix(name, "/values"); ok && name != "" && !strings.Contains(name, "/") {
			return tenancyMatch
		}
	}
	return tenancyForbidden
}

// parseTenancyRequest returns the URL parameters minus namespace, the POST
// form (empty for GET) and the namespace. Parsing is strict: url.Values drops
// pairs it cannot parse, and what it dropped would never be checked.
func parseTenancyRequest(w http.ResponseWriter, r *http.Request) (query, form url.Values, namespace string, err error) {
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, nil, "", fmt.Errorf("%w: query string: %w", errBadTenancyRequest, err)
	}

	values := query[namespaceParam]
	if len(values) != 1 {
		return nil, nil, "", fmt.Errorf("%w: exactly one namespace parameter is required", errBadTenancyRequest)
	}
	namespace = values[0]
	if problems := validation.IsDNS1123Label(namespace); len(problems) > 0 {
		return nil, nil, "", fmt.Errorf("%w: namespace %q: %s", errBadTenancyRequest, namespace, strings.Join(problems, "; "))
	}
	query.Del(namespaceParam)

	form = url.Values{}
	if r.Method == http.MethodPost {
		form, err = readForm(w, r)
		if err != nil {
			return nil, nil, "", err
		}
		if form.Has(namespaceParam) {
			return nil, nil, "", fmt.Errorf("%w: namespace belongs in the URL, not the body", errBadTenancyRequest)
		}
	}
	return query, form, namespace, nil
}

// readForm accepts only a urlencoded body. Prometheus would also read a
// multipart one (FormValue parses it), which a urlencoded parse here would
// not see.
func readForm(w http.ResponseWriter, r *http.Request) (url.Values, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != formContentType {
		return nil, fmt.Errorf("%w: a POST body must be %s", errBadTenancyRequest, formContentType)
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxFormBytes))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return nil, fmt.Errorf("%w: more than %d bytes", errBodyTooLarge, maxFormBytes)
		}
		return nil, fmt.Errorf("%w: body: %w", errBadTenancyRequest, err)
	}
	form, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("%w: body: %w", errBadTenancyRequest, err)
	}
	return form, nil
}

// forwardedForm is the body the agent receives: the hub's own encoding of the
// rewritten form, which can be longer than what the caller sent. Prometheus
// parses at most maxFormBytes of form and otherwise reads the URL alone, so a
// longer body is refused here instead of reaching it without its matchers.
func forwardedForm(form url.Values) (string, error) {
	body := form.Encode()
	if len(body) > maxFormBytes {
		return "", fmt.Errorf("%w: more than %d bytes once rewritten", errBodyTooLarge, maxFormBytes)
	}
	return body, nil
}

func enforceNamespace(kind tenancyKind, namespace string, query, form url.Values) error {
	matcher := labels.MustNewMatcher(labels.MatchEqual, namespaceLabel, namespace)
	enforcer := injectproxy.NewPromQLEnforcer(true, matcher)

	switch kind {
	case tenancyQuery:
		return enforceQuery(enforcer, query, form)
	case tenancyMatch:
		return enforceMatchers(enforcer, matcher, query, form)
	default:
		// ServeTenancy refuses tenancyForbidden before it gets here.
		return fmt.Errorf("%w: kind %d", errNoNamespaceRule, kind)
	}
}

// enforceQuery wants exactly one query between the URL and the body: Prometheus
// reads only one of them (the body's first), so a second is a query that would
// either go unchecked or be checked and ignored.
func enforceQuery(enforcer *injectproxy.PromQLEnforcer, query, form url.Values) error {
	if n := len(query[queryParam]) + len(form[queryParam]); n != 1 {
		return fmt.Errorf("%w: exactly one query parameter is required", errBadTenancyRequest)
	}
	for _, values := range []url.Values{query, form} {
		if !values.Has(queryParam) {
			continue
		}
		enforced, err := enforcer.Enforce(values.Get(queryParam))
		if err != nil {
			return err
		}
		values.Set(queryParam, enforced)
	}
	return nil
}

// enforceMatchers narrows every match[] to the namespace, refusing a selector
// that cannot match it, and supplies one when there is none — otherwise
// labels and label values would answer for the whole cluster.
func enforceMatchers(enforcer *injectproxy.PromQLEnforcer, matcher *labels.Matcher, query, form url.Values) error {
	if len(query[matchParam])+len(form[matchParam]) == 0 {
		query.Set(matchParam, selectorString([]*labels.Matcher{matcher}))
		return nil
	}

	p := parser.NewParser(parser.Options{})
	for _, values := range []url.Values{query, form} {
		for i, selector := range values[matchParam] {
			ms, err := p.ParseMetricSelector(selector)
			if err != nil {
				return fmt.Errorf("%w: %w", injectproxy.ErrQueryParse, err)
			}
			enforced, err := enforcer.EnforceMatchers(ms)
			if err != nil {
				return err
			}
			values[matchParam][i] = selectorString(enforced)
		}
	}
	return nil
}

func selectorString(ms []*labels.Matcher) string {
	parts := make([]string, len(ms))
	for i, m := range ms {
		parts[i] = m.String()
	}
	return "{" + strings.Join(parts, ",") + "}"
}
