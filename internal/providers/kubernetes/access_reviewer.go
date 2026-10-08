package kubernetes

import (
	"context"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/otterscale/otterscale/internal/core"
)

// accessReviewer implements core.AccessReviewer with a SelfSubjectAccessReview
// sent as the impersonated caller, so the target cluster's RBAC decides — the
// hub never needs to know a cluster's RoleBindings. Any authenticated user may
// create one about themselves.
type accessReviewer struct {
	kubernetes *Kubernetes
}

func NewAccessReviewer(kubernetes *Kubernetes) core.AccessReviewer {
	return &accessReviewer{kubernetes: kubernetes}
}

var _ core.AccessReviewer = (*accessReviewer)(nil)

// Allowed reports the review's verdict. A review the authorizer evaluated with
// errors but still answered "no" is a "no"; a review that could not be sent or
// answered is an error.
func (r *accessReviewer) Allowed(ctx context.Context, cluster string, attrs core.AccessAttributes) (bool, error) { //nolint:gocritic // hugeParam: signature fixed by core.AccessReviewer
	config, err := r.kubernetes.impersonationConfig(ctx, cluster)
	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, &core.DomainError{Code: core.ErrorCodeInternal, Message: "create clientset for access review", Cause: err}
	}

	review := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace:   attrs.Namespace,
				Verb:        attrs.Verb,
				Group:       attrs.Group,
				Resource:    attrs.Resource,
				Subresource: attrs.Subresource,
				Name:        attrs.Name,
			},
		},
	}

	result, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		return false, wrapK8sError(err)
	}
	return result.Status.Allowed, nil
}
