package core

import (
	"errors"
	"fmt"
)

// ErrorCode is a domain-level error category, standing in for
// infrastructure-specific codes such as K8s StatusReason. The handler layer
// maps these onto transport codes.
type ErrorCode int

const (
	ErrorCodeInternal           ErrorCode = iota // catch-all
	ErrorCodeInvalidArgument                     // bad input
	ErrorCodeNotFound                            // resource missing
	ErrorCodeAlreadyExists                       // duplicate
	ErrorCodeUnauthenticated                     // no/invalid creds
	ErrorCodePermissionDenied                    // forbidden
	ErrorCodeFailedPrecondition                  // conflict / precondition
	ErrorCodeDeadlineExceeded                    // timeout
	ErrorCodeResourceExhausted                   // rate-limit / quota
	ErrorCodeUnimplemented                       // method not allowed
	ErrorCodeUnavailable                         // service unavailable
	ErrorCodeCanceled                            // caller gave up
)

// Field names used by ErrInvalidInput. They are also the log keys the
// resource layer records the same values under, so keeping one spelling keeps
// validation errors and logs searchable together.
const (
	fieldCluster = "cluster"
	fieldName    = "name"
)

// Messages shared by more than one ErrInvalidInput.
const (
	msgResourceNameRequired = "resource name is required"
	msgPodNameRequired      = "pod name is required"
	msgMustNotBeEmpty       = "must not be empty"
)

// DomainError carries an ErrorCode and an optional cause. Infrastructure
// adapters wrap external errors into one, so the handler layer only has to
// understand domain codes.
type DomainError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error { return e.Cause }

// DomainErrorCode returns ErrorCodeInternal and false for non-domain errors.
func DomainErrorCode(err error) (ErrorCode, bool) {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Code, true
	}
	return ErrorCodeInternal, false
}

// ErrClusterNotFound means the cluster is not registered with the tunnel
// provider.
type ErrClusterNotFound struct {
	Cluster string
}

func (e *ErrClusterNotFound) Error() string {
	return fmt.Sprintf("cluster %s not registered", e.Cluster)
}

// ErrNotReady means a required subsystem is not initialized yet.
type ErrNotReady struct {
	Subsystem string
}

func (e *ErrNotReady) Error() string {
	return fmt.Sprintf("%s not initialized", e.Subsystem)
}

// ErrInvalidInput stands in for k8s apierrors.NewBadRequest, keeping core free
// of infrastructure error types.
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e *ErrInvalidInput) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("invalid %s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ErrSessionNotFound means the session is not in the store.
type ErrSessionNotFound struct {
	Resource string
	ID       string
}

func (e *ErrSessionNotFound) Error() string {
	return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}

// ErrCommandExited reports a command that ran to completion and exited
// non-zero.
//
// It is deliberately distinct from every other error an exec session can end
// with. Failing to *start* a session — no such pod, no such container, RBAC
// denial, a failed protocol upgrade — is a failure of the RPC. A command that
// started, produced its output, and exited 1 is not: the caller got exactly
// what it asked for, and reporting that as a transport failure would turn every
// unsuccessful shell command into a server error.
type ErrCommandExited struct {
	Code   int
	Reason string
}

func (e *ErrCommandExited) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("command exited with status %d: %s", e.Code, e.Reason)
	}
	return fmt.Sprintf("command exited with status %d", e.Code)
}
