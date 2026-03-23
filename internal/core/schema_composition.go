package core

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

// SchemaFieldRef identifies a GVK and a dot-separated property path
// within its OpenAPI schema. It is used on both the source and target
// sides of a SchemaCompositionRule.
type SchemaFieldRef struct {
	Group     string
	Version   string // empty = match any version (source) or required (target)
	Kind      string
	FieldPath string // dot-separated path, e.g. "spec.helmRelease.spec"
}

// SchemaCompositionRule describes how a sub-field in one GVK's schema
// should be enriched with the schema of another GVK. This is domain
// knowledge: OtterScale CRDs embed references to external CRDs
// (e.g. FluxCD HelmRelease) as RawExtension fields whose OpenAPI
// schemas are empty by default. Composition fills them at runtime.
type SchemaCompositionRule struct {
	Source SchemaFieldRef
	Target SchemaFieldRef
}

// compositionRules is the static registry of schema composition
// rules. Each entry maps an empty/RawExtension field in a source CRD
// to the corresponding schema from a target CRD.
var compositionRules = []SchemaCompositionRule{
	{
		Source: SchemaFieldRef{
			Group:     "module.otterscale.io",
			Version:   "v1alpha1",
			Kind:      "ModuleTemplate",
			FieldPath: "spec.helmRelease",
		},
		Target: SchemaFieldRef{
			Group:     "helm.toolkit.fluxcd.io",
			Version:   "v2",
			Kind:      "HelmRelease",
			FieldPath: "spec",
		},
	},
	{
		Source: SchemaFieldRef{
			Group:     "module.otterscale.io",
			Version:   "v1alpha1",
			Kind:      "ModuleTemplate",
			FieldPath: "spec.kustomization",
		},
		Target: SchemaFieldRef{
			Group:     "kustomize.toolkit.fluxcd.io",
			Version:   "v1",
			Kind:      "Kustomization",
			FieldPath: "spec",
		},
	},
}

// ComposingSchemaResolver wraps an upstream SchemaResolver and
// applies composition rules after resolution. When the resolved GVK
// matches a rule, it fetches the target GVK's schema from the same
// cluster and merges the relevant sub-schema into the base schema.
type ComposingSchemaResolver struct {
	upstream SchemaResolver
}

// NewComposingSchemaResolver returns a ComposingSchemaResolver that
// decorates the given upstream resolver with schema composition.
func NewComposingSchemaResolver(upstream SchemaResolver) *ComposingSchemaResolver {
	return &ComposingSchemaResolver{upstream: upstream}
}

var _ SchemaResolver = (*ComposingSchemaResolver)(nil)

// ResolveSchema fetches the base schema and enriches it according to
// any matching composition rules before returning.
func (r *ComposingSchemaResolver) ResolveSchema(
	ctx context.Context,
	cluster, group, version, kind string,
) (*spec.Schema, error) {
	base, err := r.upstream.ResolveSchema(ctx, cluster, group, version, kind)
	if err != nil {
		return nil, err
	}

	rules := matchingRules(group, version, kind)
	if len(rules) == 0 {
		return base, nil
	}

	// Deep-copy so we never mutate the cached schema held by the
	// upstream resolver (DiscoveryCache stores *spec.Schema values
	// that are shared across concurrent callers).
	base, err = deepCopySchema(base)
	if err != nil {
		return nil, err
	}

	for i := range rules {
		r.applyRule(ctx, cluster, &rules[i], base)
	}

	return base, nil
}

// applyRule fetches the target schema and injects the extracted
// sub-schema into the base at the rule's Source.FieldPath. Errors
// are logged and swallowed so that a missing target CRD does not
// break the primary schema request.
func (r *ComposingSchemaResolver) applyRule(
	ctx context.Context,
	cluster string,
	rule *SchemaCompositionRule,
	base *spec.Schema,
) {
	target, err := r.upstream.ResolveSchema(
		ctx, cluster,
		rule.Target.Group, rule.Target.Version, rule.Target.Kind,
	)
	if err != nil {
		slog.Warn("schema composition: failed to fetch target schema",
			"source", rule.Source.Group+"/"+rule.Source.Kind,
			"target", rule.Target.Group+"/"+rule.Target.Kind,
			"error", err,
		)
		return
	}

	sub := extractSubSchema(target, rule.Target.FieldPath)
	if sub == nil {
		slog.Warn("schema composition: target field path not found",
			"target", rule.Target.Group+"/"+rule.Target.Kind,
			"path", rule.Target.FieldPath,
		)
		return
	}

	setSubSchema(base, rule.Source.FieldPath, sub)
}

// matchingRules returns the composition rules whose source matches
// the given group, version, and kind. A rule with an empty
// Source.Version matches any version.
func matchingRules(group, version, kind string) []SchemaCompositionRule {
	var matched []SchemaCompositionRule
	for i := range compositionRules {
		if compositionRules[i].Source.Group != group || compositionRules[i].Source.Kind != kind {
			continue
		}
		if compositionRules[i].Source.Version != "" && compositionRules[i].Source.Version != version {
			continue
		}
		matched = append(matched, compositionRules[i])
	}
	return matched
}

// deepCopySchema performs a JSON round-trip to produce an independent
// copy of the schema. This is intentionally simple; the cost is
// negligible compared to the upstream network I/O that populated the
// schema in the first place.
func deepCopySchema(s *spec.Schema) (*spec.Schema, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	out := &spec.Schema{}
	return out, json.Unmarshal(data, out)
}

// extractSubSchema traverses a spec.Schema by dot-separated property
// path and returns the nested schema. Returns nil if any segment is
// missing.
func extractSubSchema(s *spec.Schema, path string) *spec.Schema {
	if path == "" {
		return s
	}

	current := s
	for seg := range strings.SplitSeq(path, ".") {
		if current.Properties == nil {
			return nil
		}
		next, ok := current.Properties[seg]
		if !ok {
			return nil
		}
		current = &next
	}
	return current
}

// setSubSchema sets a nested schema value at the given dot-separated
// path. It uses recursive write-back to handle spec.SchemaProperties'
// value semantics (map[string]spec.Schema stores copies, not
// pointers).
func setSubSchema(base *spec.Schema, path string, sub *spec.Schema) {
	if path == "" || sub == nil {
		return
	}
	setSubSchemaAt(base, strings.Split(path, "."), sub)
}

func setSubSchemaAt(current *spec.Schema, segments []string, sub *spec.Schema) {
	if len(segments) == 0 || current.Properties == nil {
		return
	}

	seg := segments[0]
	prop, ok := current.Properties[seg]
	if !ok {
		return
	}

	if len(segments) == 1 {
		current.Properties[seg] = *sub
		return
	}

	setSubSchemaAt(&prop, segments[1:], sub)
	current.Properties[seg] = prop
}
