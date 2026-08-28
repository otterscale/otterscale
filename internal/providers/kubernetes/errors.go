package kubernetes

import (
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// statusReasonToDomainCode keeps the Kubernetes-specific mapping in the adapter
// layer, out of the handler and core layers.
var statusReasonToDomainCode = map[metav1.StatusReason]core.ErrorCode{
	metav1.StatusReasonUnauthorized:          core.ErrorCodeUnauthenticated,
	metav1.StatusReasonForbidden:             core.ErrorCodePermissionDenied,
	metav1.StatusReasonNotFound:              core.ErrorCodeNotFound,
	metav1.StatusReasonAlreadyExists:         core.ErrorCodeAlreadyExists,
	metav1.StatusReasonConflict:              core.ErrorCodeFailedPrecondition,
	metav1.StatusReasonGone:                  core.ErrorCodeNotFound,
	metav1.StatusReasonInvalid:               core.ErrorCodeInvalidArgument,
	metav1.StatusReasonServerTimeout:         core.ErrorCodeDeadlineExceeded,
	metav1.StatusReasonStoreReadError:        core.ErrorCodeInternal,
	metav1.StatusReasonTimeout:               core.ErrorCodeDeadlineExceeded,
	metav1.StatusReasonTooManyRequests:       core.ErrorCodeResourceExhausted,
	metav1.StatusReasonBadRequest:            core.ErrorCodeInvalidArgument,
	metav1.StatusReasonMethodNotAllowed:      core.ErrorCodeUnimplemented,
	metav1.StatusReasonNotAcceptable:         core.ErrorCodeInvalidArgument,
	metav1.StatusReasonRequestEntityTooLarge: core.ErrorCodeResourceExhausted,
	metav1.StatusReasonUnsupportedMediaType:  core.ErrorCodeInvalidArgument,
	metav1.StatusReasonInternalError:         core.ErrorCodeInternal,
	metav1.StatusReasonExpired:               core.ErrorCodeInvalidArgument,
	metav1.StatusReasonServiceUnavailable:    core.ErrorCodeUnavailable,
}

// wrapK8sError returns non-Kubernetes errors as-is, so callers should only
// pass errors originating from API calls.
func wrapK8sError(err error) error {
	if err == nil {
		return nil
	}

	var apiStatus apierrors.APIStatus
	if !errors.As(err, &apiStatus) {
		return err
	}

	code, ok := statusReasonToDomainCode[apiStatus.Status().Reason]
	if !ok {
		code = core.ErrorCodeInternal
	}

	return &core.DomainError{
		Code:    code,
		Message: apiStatus.Status().Message,
		Cause:   err,
	}
}
