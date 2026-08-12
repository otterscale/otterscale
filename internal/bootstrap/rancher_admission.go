package bootstrap

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	rancherWebhookRoleName           = "otterscale-tenant-operator-rancher-webhook"
	rancherProjectPolicyName         = "otterscale-workspace-rancher-project"
	rancherNamespacePolicyName       = "otterscale-workspace-namespace-security"
	rancherProjectSnapshotAnnotation = "security.otterscale.io/rancher-project-id"
	tenantOperatorServiceAccount     = "tenant-operator-controller-manager"
	tenantOperatorServiceAccountNS   = "otterscale-system"
	rancherGuardDefaultPollInterval  = time.Second
	rancherGuardDefaultPollTimeout   = time.Minute
)

var (
	validatingAdmissionPolicyGVR = schema.GroupVersionResource{
		Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingadmissionpolicies",
	}
	validatingAdmissionPolicyBindingGVR = schema.GroupVersionResource{
		Group: "admissionregistration.k8s.io", Version: "v1", Resource: "validatingadmissionpolicybindings",
	}
	clusterRoleBindingGVR = schema.GroupVersionResource{
		Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings",
	}
	clusterRoleGVR = schema.GroupVersionResource{
		Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles",
	}
)

// activateRancherWebhookAccess removes any previous activation, waits until
// both fail-closed guards match this agent's Project snapshot, and only then
// binds the tenant-operator to the synthetic Rancher webhook role.
func (b *Bootstrapper) activateRancherWebhookAccess(ctx context.Context, rancherProjectID string) error {
	if err := b.deactivateRancherWebhookAccess(ctx); err != nil {
		return err
	}
	bindings := b.dynamic.Resource(clusterRoleBindingGVR)

	interval, timeout := b.rancherGuardPollInterval, b.rancherGuardPollTimeout
	if interval <= 0 {
		interval = rancherGuardDefaultPollInterval
	}
	if timeout <= 0 {
		timeout = rancherGuardDefaultPollTimeout
	}

	var lastErr error
	err := wait.PollUntilContextTimeout(ctx, interval, timeout, true, func(ctx context.Context) (bool, error) {
		lastErr = b.validateRancherAdmissionGuard(ctx, rancherProjectID)
		return lastErr == nil, nil
	})
	if err != nil {
		if lastErr == nil {
			lastErr = err
		}
		return fmt.Errorf("Rancher admission guard not ready: %w", lastErr)
	}

	binding := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "ClusterRoleBinding",
		"metadata": map[string]any{
			"name": rancherWebhookRoleName,
		},
		"subjects": []any{map[string]any{
			"kind":      "ServiceAccount",
			"name":      tenantOperatorServiceAccount,
			"namespace": tenantOperatorServiceAccountNS,
		}},
		"roleRef": map[string]any{
			"apiGroup": "rbac.authorization.k8s.io",
			"kind":     "ClusterRole",
			"name":     rancherWebhookRoleName,
		},
	}}
	if _, err := bindings.Create(ctx, binding, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("activate Rancher webhook access: %w", err)
	}
	return nil
}

func (b *Bootstrapper) deactivateRancherWebhookAccess(ctx context.Context) error {
	err := b.dynamic.Resource(clusterRoleBindingGVR).Delete(ctx, rancherWebhookRoleName, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("deactivate Rancher webhook access: %w", err)
	}
	return nil
}

func (b *Bootstrapper) validateRancherAdmissionGuard(ctx context.Context, rancherProjectID string) error {
	role, err := b.dynamic.Resource(clusterRoleGVR).Get(ctx, rancherWebhookRoleName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get Rancher webhook ClusterRole: %w", err)
	}
	roleSnapshot, found := role.GetAnnotations()[rancherProjectSnapshotAnnotation]
	if !found || roleSnapshot != rancherProjectID {
		return fmt.Errorf("Rancher webhook ClusterRole has Project snapshot %q, want %q", roleSnapshot, rancherProjectID)
	}
	if err := validateRancherWebhookRole(role, rancherProjectID); err != nil {
		return err
	}

	for _, name := range []string{rancherProjectPolicyName, rancherNamespacePolicyName} {
		policy, err := b.dynamic.Resource(validatingAdmissionPolicyGVR).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get ValidatingAdmissionPolicy %s: %w", name, err)
		}
		got, found := policy.GetAnnotations()[rancherProjectSnapshotAnnotation]
		if !found || got != rancherProjectID {
			return fmt.Errorf("ValidatingAdmissionPolicy %s has Project snapshot %q, want %q", name, got, rancherProjectID)
		}
		failurePolicy, found, err := unstructured.NestedString(policy.Object, "spec", "failurePolicy")
		if err != nil || !found || failurePolicy != "Fail" {
			return fmt.Errorf("ValidatingAdmissionPolicy %s is not fail closed", name)
		}
		observed, found, err := unstructured.NestedInt64(policy.Object, "status", "observedGeneration")
		if err != nil || !found || observed != policy.GetGeneration() {
			return fmt.Errorf("ValidatingAdmissionPolicy %s generation is not observed", name)
		}
		if _, found, err := unstructured.NestedMap(policy.Object, "status", "typeChecking"); err != nil || !found {
			return fmt.Errorf("ValidatingAdmissionPolicy %s type checking is not complete", name)
		}
		warnings, found, err := unstructured.NestedSlice(policy.Object, "status", "typeChecking", "expressionWarnings")
		if err != nil || (found && len(warnings) != 0) {
			return fmt.Errorf("ValidatingAdmissionPolicy %s has type-check warnings", name)
		}

		binding, err := b.dynamic.Resource(validatingAdmissionPolicyBindingGVR).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get ValidatingAdmissionPolicyBinding %s: %w", name, err)
		}
		policyName, found, err := unstructured.NestedString(binding.Object, "spec", "policyName")
		if err != nil || !found || policyName != name {
			return fmt.Errorf("ValidatingAdmissionPolicyBinding %s references the wrong policy", name)
		}
		actions, found, err := unstructured.NestedStringSlice(binding.Object, "spec", "validationActions")
		if err != nil || !found || len(actions) != 1 || actions[0] != "Deny" {
			return fmt.Errorf("ValidatingAdmissionPolicyBinding %s is not Deny-only", name)
		}
		if _, found, err := unstructured.NestedFieldNoCopy(binding.Object, "spec", "matchResources"); err != nil || found {
			return fmt.Errorf("ValidatingAdmissionPolicyBinding %s narrows guarded resources", name)
		}
	}
	return nil
}

func validateRancherWebhookRole(role *unstructured.Unstructured, rancherProjectID string) error {
	rules, found, err := unstructured.NestedSlice(role.Object, "rules")
	if err != nil || !found || len(rules) != 1 {
		return fmt.Errorf("Rancher webhook ClusterRole must contain exactly one rule")
	}
	rule, ok := rules[0].(map[string]any)
	if !ok || !exactStringSlice(rule["apiGroups"], "management.cattle.io") || !exactStringSlice(rule["resources"], "projects") {
		return fmt.Errorf("Rancher webhook ClusterRole has unexpected API resources")
	}

	if rancherProjectID == "" {
		if _, found := rule["resourceNames"]; found || !exactStringSlice(rule["verbs"], "updatepsa") {
			return fmt.Errorf("empty-Project Rancher webhook ClusterRole is broader than allowed")
		}
		return nil
	}

	_, projectName, err := core.ParseRancherProjectID(rancherProjectID)
	if err != nil {
		return fmt.Errorf("parse Rancher Project ID for guard: %w", err)
	}
	if !exactStringSlice(rule["resourceNames"], projectName) || !exactStringSlice(rule["verbs"], "updatepsa", "manage-namespaces") {
		return fmt.Errorf("selected-Project Rancher webhook ClusterRole does not match Project %q", projectName)
	}
	return nil
}

func exactStringSlice(value any, want ...string) bool {
	items, ok := value.([]any)
	if !ok || len(items) != len(want) {
		return false
	}
	for i := range want {
		if items[i] != want[i] {
			return false
		}
	}
	return true
}
