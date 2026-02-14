package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// applyManifest parses a multi-document YAML byte slice and applies
// every object to the cluster via Server-Side Apply. CRDs are applied
// first and the function blocks until each CRD reaches the
// Established condition, ensuring that subsequent resources whose GVR
// depends on those CRDs can be resolved.
func (b *Bootstrapper) applyManifest(ctx context.Context, data []byte) error {
	objects, err := parseMultiDoc(data)
	if err != nil {
		return fmt.Errorf("parse multi-doc YAML: %w", err)
	}

	if len(objects) == 0 {
		return nil
	}

	// Partition into CRDs and non-CRD resources.
	var crds, rest []*unstructured.Unstructured
	for _, obj := range objects {
		if obj.GetKind() == "CustomResourceDefinition" {
			crds = append(crds, obj)
		} else {
			rest = append(rest, obj)
		}
	}

	// Phase 1: Apply CRDs and wait for them to be established.
	if len(crds) > 0 {
		mapper := b.newMapper()
		for _, crd := range crds {
			if err := b.applyObject(ctx, mapper, crd); err != nil {
				return fmt.Errorf("apply CRD %s: %w", crd.GetName(), err)
			}
			b.log.Info("applied CRD", "name", crd.GetName())
		}

		if err := b.waitForCRDs(ctx, crds); err != nil {
			return err
		}
	}

	// Phase 2: Apply remaining resources with a fresh mapper that
	// knows about the newly established CRDs.
	if len(rest) > 0 {
		mapper := b.newMapper()
		for _, obj := range rest {
			if err := b.applyObject(ctx, mapper, obj); err != nil {
				return fmt.Errorf("apply %s %s/%s: %w",
					obj.GetKind(), obj.GetNamespace(), obj.GetName(), err)
			}
			b.log.Info("applied resource",
				"kind", obj.GetKind(),
				"namespace", obj.GetNamespace(),
				"name", obj.GetName(),
			)
		}
	}

	return nil
}

// applyObject performs a Server-Side Apply for a single unstructured
// object. It uses the REST mapper to resolve the GVK into a GVR and
// then issues a PATCH with ApplyPatchType.
func (b *Bootstrapper) applyObject(
	ctx context.Context,
	mapper meta.RESTMapper,
	obj *unstructured.Unstructured,
) error {
	gvk := obj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("map GVK %s: %w", gvk, err)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshal object: %w", err)
	}

	force := true
	patchOpts := metav1.PatchOptions{
		FieldManager: fieldManager,
		Force:        &force,
	}

	var client dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		client = b.dynamic.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		client = b.dynamic.Resource(mapping.Resource)
	}

	_, err = client.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, patchOpts)
	return err
}

// waitForCRDs blocks until every CRD in the slice has the
// Established condition set to True. It polls with a 2-second
// interval and gives up after 60 seconds.
func (b *Bootstrapper) waitForCRDs(ctx context.Context, crds []*unstructured.Unstructured) error {
	for _, crd := range crds {
		name := crd.GetName()
		b.log.Info("waiting for CRD to be established", "name", name)

		err := wait.PollUntilContextTimeout(ctx, 2*time.Second, 60*time.Second, true,
			func(ctx context.Context) (bool, error) {
				obj, err := b.dynamic.Resource(crdGVR).Get(ctx, name, metav1.GetOptions{})
				if err != nil {
					return false, nil // retry on transient errors
				}
				return isCRDEstablished(obj), nil
			},
		)
		if err != nil {
			return fmt.Errorf("CRD %s did not become established: %w", name, err)
		}
		b.log.Info("CRD established", "name", name)
	}
	return nil
}

// isCRDEstablished inspects the CRD status conditions for
// type=Established, status=True.
func isCRDEstablished(obj *unstructured.Unstructured) bool {
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil || !found {
		return false
	}
	for _, c := range conditions {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if m["type"] == "Established" && m["status"] == "True" {
			return true
		}
	}
	return false
}

// newMapper creates a fresh REST mapper backed by a cached discovery
// client. Callers should create a new mapper after applying CRDs so
// that newly registered API resources are visible.
func (b *Bootstrapper) newMapper() meta.RESTMapper {
	cachedDisc := memory.NewMemCacheClient(b.disc)
	return restmapper.NewDeferredDiscoveryRESTMapper(cachedDisc)
}

// parseMultiDoc splits a multi-document YAML byte slice into
// individual unstructured objects, skipping empty documents.
func parseMultiDoc(data []byte) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 4096)
	for {
		obj := &unstructured.Unstructured{}
		if err := decoder.Decode(obj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// Skip empty documents (e.g. trailing "---").
		if obj.GetKind() == "" {
			continue
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// crdGVR is the GroupVersionResource for apiextensions.k8s.io/v1
// CustomResourceDefinitions, used to poll CRD status.
var crdGVR = schema.GroupVersionResource{
	Group:    "apiextensions.k8s.io",
	Version:  "v1",
	Resource: "customresourcedefinitions",
}
