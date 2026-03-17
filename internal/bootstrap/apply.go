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

// crdPollInterval is how often to check CRD establishment status.
const crdPollInterval = 2 * time.Second

// crdPollTimeout is the maximum time to wait for a CRD to become established.
const crdPollTimeout = 60 * time.Second

// waitForCRDs blocks until every CRD in the slice has the
// Established condition set to True. It polls with a 2-second
// interval and gives up after 60 seconds.
func (b *Bootstrapper) waitForCRDs(ctx context.Context, crds []*unstructured.Unstructured) error {
	for _, crd := range crds {
		name := crd.GetName()
		b.log.Info("waiting for CRD to be established", "name", name)

		err := wait.PollUntilContextTimeout(ctx, crdPollInterval, crdPollTimeout, true,
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
		b.log.Info("crd established", "name", name)
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

// newMapper invalidates the shared cached discovery client and returns
// a fresh REST mapper so that newly registered API resources (e.g.
// CRDs applied in an earlier phase) are visible.
func (b *Bootstrapper) newMapper() meta.RESTMapper {
	b.cachedDisc.Invalidate()
	return restmapper.NewDeferredDiscoveryRESTMapper(b.cachedDisc)
}

// parseMultiDoc splits a multi-document YAML byte slice into
// individual unstructured objects, skipping empty documents.
func parseMultiDoc(data []byte) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	const yamlDecoderBufSize = 4096
	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), yamlDecoderBufSize)
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

// deploymentGVR is the GroupVersionResource for apps/v1 Deployments,
// used to poll Deployment availability status.
var deploymentGVR = schema.GroupVersionResource{
	Group:    "apps",
	Version:  "v1",
	Resource: "deployments",
}

// deploymentPollInterval is how often to check Deployment availability.
const deploymentPollInterval = 3 * time.Second

// deploymentPollTimeout is the maximum time to wait for a Deployment
// to become available. cert-manager may pull images and create
// resources, so we allow up to 5 minutes.
const deploymentPollTimeout = 300 * time.Second

// waitForDeployment blocks until the specified Deployment has the
// Available condition set to True. It polls with a 3-second interval
// and gives up after 5 minutes.
func (b *Bootstrapper) waitForDeployment(ctx context.Context, namespace, name string) error {
	b.log.Info("waiting for Deployment to be available",
		"namespace", namespace, "name", name)

	err := wait.PollUntilContextTimeout(ctx, deploymentPollInterval, deploymentPollTimeout, true,
		func(ctx context.Context) (bool, error) {
			obj, err := b.dynamic.Resource(deploymentGVR).
				Namespace(namespace).
				Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return false, nil // retry on transient errors
			}
			return isDeploymentAvailable(obj), nil
		},
	)
	if err != nil {
		return fmt.Errorf("deployment %s/%s did not become available: %w", namespace, name, err)
	}

	b.log.Info("deployment is available", "namespace", namespace, "name", name)
	return nil
}

// isDeploymentAvailable inspects the Deployment status conditions for
// type=Available, status=True.
func isDeploymentAvailable(obj *unstructured.Unstructured) bool {
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil || !found {
		return false
	}
	for _, c := range conditions {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if m["type"] == "Available" && m["status"] == "True" {
			return true
		}
	}
	return false
}
