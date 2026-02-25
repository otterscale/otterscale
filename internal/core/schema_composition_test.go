package core

import (
	"context"
	"errors"
	"testing"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

// stubSchemaResolver returns pre-configured schemas keyed by
// "group/kind". It records calls for assertion.
type stubSchemaResolver struct {
	schemas map[string]*spec.Schema
	calls   []string
	err     error
}

func (s *stubSchemaResolver) ResolveSchema(_ context.Context, _, group, _, kind string) (*spec.Schema, error) {
	key := group + "/" + kind
	s.calls = append(s.calls, key)
	if s.err != nil {
		return nil, s.err
	}
	schema, ok := s.schemas[key]
	if !ok {
		return nil, errors.New("not found: " + key)
	}
	return schema, nil
}

func newSchemaWithProps(props map[string]spec.Schema) *spec.Schema {
	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       spec.StringOrArray{"object"},
			Properties: props,
		},
	}
}

// ---------------------------------------------------------------------------
// ComposingSchemaResolver tests
// ---------------------------------------------------------------------------

func TestComposingSchemaResolver_NonMatchingGVK(t *testing.T) {
	upstream := &stubSchemaResolver{
		schemas: map[string]*spec.Schema{
			"apps/Deployment": newSchemaWithProps(map[string]spec.Schema{
				"spec": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"object"}}},
			}),
		},
	}

	resolver := NewComposingSchemaResolver(upstream)
	schema, err := resolver.ResolveSchema(context.Background(), "cluster1", "apps", "v1", "Deployment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(upstream.calls) != 1 {
		t.Errorf("expected 1 upstream call, got %d", len(upstream.calls))
	}
	if _, ok := schema.Properties["spec"]; !ok {
		t.Error("expected spec property in returned schema")
	}
}

func TestComposingSchemaResolver_MatchingGVK(t *testing.T) {
	baseSchema := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"helmRelease": {}, // empty — RawExtension placeholder
				},
			},
		},
	})

	helmReleaseSchema := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"chart": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"object"},
						},
					},
					"interval": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"string"},
						},
					},
				},
			},
		},
	})

	upstream := &stubSchemaResolver{
		schemas: map[string]*spec.Schema{
			"addons.otterscale.io/ModuleTemplate": baseSchema,
			"helm.toolkit.fluxcd.io/HelmRelease":  helmReleaseSchema,
		},
	}

	resolver := NewComposingSchemaResolver(upstream)
	schema, err := resolver.ResolveSchema(
		context.Background(), "cluster1",
		"addons.otterscale.io", "v1alpha1", "ModuleTemplate",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Three upstream calls: base + HelmRelease target + Kustomization target (not found, gracefully skipped).
	if len(upstream.calls) != 3 {
		t.Errorf("expected 3 upstream calls, got %d: %v", len(upstream.calls), upstream.calls)
	}

	// Verify spec.helmRelease is now replaced with HelmRelease's spec properties.
	enriched := navigateSchema(t, schema, "spec", "helmRelease")
	if enriched == nil {
		t.Fatal("spec.helmRelease should exist after composition")
	}
	if _, ok := enriched.Properties["chart"]; !ok {
		t.Error("expected 'chart' property from HelmRelease spec")
	}
	if _, ok := enriched.Properties["interval"]; !ok {
		t.Error("expected 'interval' property from HelmRelease spec")
	}
}

func TestComposingSchemaResolver_TargetFetchError(t *testing.T) {
	baseSchema := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"helmRelease": {}, // empty placeholder
				},
			},
		},
	})

	upstream := &stubSchemaResolver{
		schemas: map[string]*spec.Schema{
			"addons.otterscale.io/ModuleTemplate": baseSchema,
		},
	}

	resolver := NewComposingSchemaResolver(upstream)
	schema, err := resolver.ResolveSchema(
		context.Background(), "cluster1",
		"addons.otterscale.io", "v1alpha1", "ModuleTemplate",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The base schema is returned unmodified.
	enriched := navigateSchema(t, schema, "spec", "helmRelease")
	if enriched == nil {
		t.Fatal("spec.helmRelease should still exist")
	}
	if len(enriched.Properties) != 0 {
		t.Errorf("expected empty properties (no enrichment), got %d", len(enriched.Properties))
	}
}

func TestComposingSchemaResolver_CacheNotCorrupted(t *testing.T) {
	original := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"helmRelease": {}, // empty placeholder
				},
			},
		},
	})

	helmReleaseSchema := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"chart": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"object"}}},
				},
			},
		},
	})

	upstream := &stubSchemaResolver{
		schemas: map[string]*spec.Schema{
			"addons.otterscale.io/ModuleTemplate": original,
			"helm.toolkit.fluxcd.io/HelmRelease":  helmReleaseSchema,
		},
	}

	resolver := NewComposingSchemaResolver(upstream)
	_, err := resolver.ResolveSchema(
		context.Background(), "cluster1",
		"addons.otterscale.io", "v1alpha1", "ModuleTemplate",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The original schema pointer held by the upstream must NOT have
	// been mutated by the composing resolver.
	origHR := navigateSchema(t, original, "spec", "helmRelease")
	if origHR == nil {
		t.Fatal("original spec.helmRelease should still exist")
	}
	if len(origHR.Properties) != 0 {
		t.Errorf("original schema was mutated: spec.helmRelease has %d properties", len(origHR.Properties))
	}
}

func TestComposingSchemaResolver_SourceVersionMismatch(t *testing.T) {
	// Temporarily override compositionRules with a version-pinned rule.
	orig := compositionRules
	compositionRules = []SchemaCompositionRule{
		{
			Source: SchemaFieldRef{
				Group:     "addons.otterscale.io",
				Version:   "v1beta1",
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
	}
	t.Cleanup(func() { compositionRules = orig })

	baseSchema := newSchemaWithProps(map[string]spec.Schema{
		"spec": {
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"helmRelease": {}, // empty placeholder
				},
			},
		},
	})

	upstream := &stubSchemaResolver{
		schemas: map[string]*spec.Schema{
			"addons.otterscale.io/ModuleTemplate": baseSchema,
		},
	}

	resolver := NewComposingSchemaResolver(upstream)

	// Request with v1alpha1 should NOT match the v1beta1-pinned rule.
	schema, err := resolver.ResolveSchema(
		context.Background(), "cluster1",
		"addons.otterscale.io", "v1alpha1", "ModuleTemplate",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only one upstream call (base); no target fetch attempted.
	if len(upstream.calls) != 1 {
		t.Errorf("expected 1 upstream call, got %d: %v", len(upstream.calls), upstream.calls)
	}

	enriched := navigateSchema(t, schema, "spec", "helmRelease")
	if enriched == nil {
		t.Fatal("spec.helmRelease should still exist")
	}
	if len(enriched.Properties) != 0 {
		t.Errorf("expected no enrichment for version mismatch, got %d properties", len(enriched.Properties))
	}
}

// ---------------------------------------------------------------------------
// extractSubSchema / setSubSchema tests
// ---------------------------------------------------------------------------

func TestExtractSubSchema_EmptyPath(t *testing.T) {
	s := newSchemaWithProps(map[string]spec.Schema{
		"foo": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
	})
	result := extractSubSchema(s, "")
	if result != s {
		t.Error("empty path should return the schema itself")
	}
}

func TestExtractSubSchema_DeepPath(t *testing.T) {
	s := newSchemaWithProps(map[string]spec.Schema{
		"a": {
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"b": {
						SchemaProps: spec.SchemaProps{
							Properties: map[string]spec.Schema{
								"c": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
							},
						},
					},
				},
			},
		},
	})

	result := extractSubSchema(s, "a.b.c")
	if result == nil {
		t.Fatal("expected non-nil result for path a.b.c")
	}
	if len(result.Type) == 0 || result.Type[0] != "integer" {
		t.Errorf("expected type integer, got %v", result.Type)
	}
}

func TestExtractSubSchema_MissingSegment(t *testing.T) {
	s := newSchemaWithProps(map[string]spec.Schema{
		"a": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"object"}}},
	})
	result := extractSubSchema(s, "a.missing")
	if result != nil {
		t.Error("expected nil for missing path segment")
	}
}

func TestSetSubSchema_WriteBack(t *testing.T) {
	base := newSchemaWithProps(map[string]spec.Schema{
		"a": {
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"b": {
						SchemaProps: spec.SchemaProps{
							Properties: map[string]spec.Schema{
								"c": {},
							},
						},
					},
				},
			},
		},
	})

	replacement := &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"string"},
			Properties: map[string]spec.Schema{
				"injected": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"boolean"}}},
			},
		},
	}

	setSubSchema(base, "a.b.c", replacement)

	result := navigateSchema(t, base, "a", "b", "c")
	if result == nil {
		t.Fatal("expected a.b.c to exist after setSubSchema")
	}
	if _, ok := result.Properties["injected"]; !ok {
		t.Error("expected 'injected' property after setSubSchema")
	}
}

func TestSetSubSchema_MissingIntermediatePath(t *testing.T) {
	base := newSchemaWithProps(map[string]spec.Schema{
		"x": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
	})

	replacement := &spec.Schema{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"number"}}}
	setSubSchema(base, "x.y.z", replacement)

	// x has no properties, so the path is unreachable. Schema is unchanged.
	if base.Properties["x"].Type[0] != "string" {
		t.Error("schema should be unchanged when intermediate path is missing")
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func navigateSchema(t *testing.T, s *spec.Schema, segments ...string) *spec.Schema {
	t.Helper()
	current := s
	for _, seg := range segments {
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
